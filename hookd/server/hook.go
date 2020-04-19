package server

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/grafov/m3u8"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"go.uber.org/zap"
)

var (
	ErrBadRequest     = echo.NewHTTPError(http.StatusBadRequest)
	ErrInternalServer = echo.NewHTTPError(http.StatusInternalServerError)
)

type HookConfig struct {
	Prefix string
}

type Hook struct {
	cfg           *HookConfig
	logger        *zap.Logger
	e             *echo.Echo
	sc            *clientv1.ServiceClient
	segmentsCount sync.Map
}

func NewHook(ctx context.Context, e *echo.Echo, cfg *HookConfig, sc *clientv1.ServiceClient) (*Hook, error) {
	hook := &Hook{
		e:      e,
		cfg:    cfg,
		logger: ctxzap.Extract(ctx).With(zap.String("system", "server")),
		sc:     sc,
	}
	hook.e.Any(cfg.Prefix, hook.handleHook)
	return hook, nil
}

func (h *Hook) handleHook(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()

	err := req.ParseForm()
	if err != nil {
		return ErrBadRequest
	}

	logger := h.logger
	for k, v := range req.Form {
		logger = logger.With(zap.String(fmt.Sprintf("form_%s", k), v[0]))
	}
	logger.Info("hook request")

	call := req.FormValue("call")
	switch call {
	case "publish":
		err = h.handlePublish(ctx, req)
	case "publish_done":
		err = h.handlePublishDone(ctx, req)
	case "playlist":
		err = h.handlePlaylist(ctx, req)
	case "update_publish":
		err = h.handleUpdatePublish(ctx, req)
	}

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Hook) handlePublish(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hook.handlePublish")
	defer span.Finish()

	logger := ctxzap.Extract(ctx)
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warn("failed to get stream id")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	logger = logger.With(zap.String("stream_id", streamID))

	req := &privatev1.StreamRequest{Id: streamID}
	streamResp, err := h.sc.Streams.Get(ctx, req)
	if err != nil {
		logger.Error("failed to get stream", zap.Error(err))
		return ErrBadRequest
	}

	if streamResp.Status != v1.StreamStatusPrepared {
		logger.Warn("wrong stream status", zap.String("status", streamResp.Status.String()))
		return ErrBadRequest
	}

	_, err = h.sc.Streams.Publish(ctx, req)
	if err != nil {
		logger.Error("failed to publish", zap.Error(err))
		return ErrBadRequest
	}

	return nil
}

func (h *Hook) handlePublishDone(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hook.handlePublishDone")
	defer span.Finish()

	logger := ctxzap.Extract(ctx)
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warn("failed to get stream id")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	logger = logger.With(zap.String("stream_id", streamID))

	logger.Info("publishing done")

	req := &privatev1.StreamRequest{Id: streamID}
	_, err := h.sc.Streams.Stop(ctx, req)
	if err != nil {
		logger.Error("failed to stop stream", zap.Error(err))
		return ErrBadRequest
	}

	return nil
}

func (h *Hook) handlePlaylist(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hook.handlePlaylist")
	defer span.Finish()

	logger := ctxzap.Extract(ctx)
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warn("failed to get stream name")
		return ErrBadRequest
	}

	path := r.FormValue("path")
	if path == "" {
		logger.Warn("failed to get stream path")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	span.SetTag("path", path)

	logger = logger.With(zap.String("stream_id", streamID), zap.String("path", path))
	logger.Info("updating playlist")

	f, err := os.Open(path)
	if err != nil {
		logger.Error("failed to open playlist", zap.Error(err))
		return ErrInternalServer
	}

	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		logger.Error("failed to decode playlist", zap.Error(err))
		return ErrInternalServer
	}

	switch listType {
	case m3u8.MASTER:
		logger.Error("failed to playlist type")
		return ErrInternalServer
	case m3u8.MEDIA:
		pl := p.(*m3u8.MediaPlaylist)
		segmentsCount := 0
		for _, s := range pl.Segments {
			if s != nil {
				segmentsCount++
			}
		}

		req := &privatev1.StreamRequest{Id: streamID}
		streamResp, err := h.sc.Streams.Get(ctx, req)
		if err != nil {
			logger.Error("failed to get stream", zap.Error(err))
			return ErrBadRequest
		}

		actual, ok := h.segmentsCount.LoadOrStore(streamID, segmentsCount)
		if ok {
			h.segmentsCount.Store(streamID, segmentsCount)
		}
		prevSegmentsCount := actual.(int)

		for i := prevSegmentsCount; i < segmentsCount; i++ {
			achReq := &emitterv1.AddInputChunkRequest{
				StreamContractId: streamResp.StreamContractID,
				ChunkId:          uint64(i),
				Reward:           streamResp.ProfileCost / 60 * pl.Segments[i-1].Duration,
			}

			_, err := h.sc.Emitter.AddInputChunk(ctx, achReq)
			if err != nil {
				logger.Error("failed to add input chunk", zap.Error(err))
			}
		}
	}

	return nil
}

func (h *Hook) handleUpdatePublish(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "hook.handleUpdatePublish")
	defer span.Finish()

	logger := ctxzap.Extract(ctx)
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warn("failed to get stream id")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	logger = logger.With(zap.String("stream_id", streamID))

	req := &privatev1.StreamRequest{Id: streamID}
	streamResp, err := h.sc.Streams.Get(ctx, req)
	if err != nil {
		logger.Error("failed to get stream", zap.Error(err))
		return nil
	}

	logger.Info("stream status is", zap.String("status", streamResp.Status.String()))

	if streamResp.Status == v1.StreamStatusFailed ||
		streamResp.Status == v1.StreamStatusCancelled ||
		streamResp.Status == v1.StreamStatusCompleted {
		return ErrBadRequest
	}

	return nil
}
