package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	dispatcherv1 "github.com/videocoin/cloud-api/dispatcher/v1"
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
	cfg                 *HookConfig
	logger              *zap.Logger
	e                   *echo.Echo
	sc                  *clientv1.ServiceClient
	segmentsCount       sync.Map
	addInputChunkFailed sync.Map
	playlists           sync.Map
}

func NewHook(ctx context.Context, e *echo.Echo, cfg *HookConfig, sc *clientv1.ServiceClient) (*Hook, error) {
	hook := &Hook{
		e:      e,
		cfg:    cfg,
		logger: ctxzap.Extract(ctx),
		sc:     sc,
	}
	hook.e.Any(cfg.Prefix, hook.handleHook)
	return hook, nil
}

func (h *Hook) handleHook(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()

	span, spanCtx := opentracing.StartSpanFromContext(ctx, "hook.handleHook")
	defer span.Finish()

	err := req.ParseForm()
	if err != nil {
		h.logger.Warn("failed to parse form", zap.Error(err))
		return ErrBadRequest
	}

	logger := h.logger
	for k, v := range req.Form {
		logger = logger.With(zap.String(fmt.Sprintf("form_%s", k), v[0]))
	}
	logger.Info("hook request")

	call := req.FormValue("call")
	streamID := req.FormValue("name")
	if streamID == "" {
		logger.Error("failed to get stream name")
		return ErrBadRequest
	}

	span.SetTag("hook", call)
	span.SetTag("stream_id", streamID)

	logger = logger.With(
		zap.String("stream_id", streamID),
		zap.String("call", call),
	)

	hookCtx := ctxzap.ToContext(spanCtx, logger)
	hookCtx = opentracing.ContextWithSpan(hookCtx, span)

	switch call {
	case "publish":
		err = h.handlePublish(hookCtx, streamID)
	case "publish_done":
		err = h.handlePublishDone(hookCtx, streamID)
	case "playlist":
		err = h.handlePlaylist(hookCtx, streamID, req)
	case "update_publish":
		err = h.handleUpdatePublish(hookCtx, streamID)
	}

	if err != nil {
		if strings.HasPrefix(err.Error(), "stream status is") {
			return ErrBadRequest
		}
		logger.Error(err.Error())
		return ErrBadRequest
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Hook) handlePublish(ctx context.Context, streamID string) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "hook.handlePublish")
	defer span.Finish()

	logger := ctxzap.Extract(ctx).With(zap.String("stream_id", streamID))
	logger.Info("publishing")

	stream, err := h.sc.Streams.Get(ctx, newStreamRequest(streamID))
	if err != nil {
		return fmt.Errorf("failed to get stream: %s", err)
	}

	if stream.Status != v1.StreamStatusPrepared {
		return fmt.Errorf("wrong stream status: %s", stream.Status.String())
	}

	_, err = h.sc.Streams.Publish(spanCtx, newStreamRequest(streamID))
	if err != nil {
		return fmt.Errorf("failed to stream publish: %s", err)
	}

	h.segmentsCount.Store(streamID, uint64(0))

	return nil
}

func (h *Hook) handlePublishDone(ctx context.Context, streamID string) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "hook.handlePublishDone")
	defer span.Finish()

	logger := ctxzap.Extract(ctx).With(zap.String("stream_id", streamID))
	logger.Info("publishing done")

	_, err := h.sc.Streams.Stop(spanCtx, newStreamRequest(streamID))
	if err != nil {
		return fmt.Errorf("failed to stop stream: %s", err)
	}

	return nil
}

func (h *Hook) handlePlaylist(ctx context.Context, streamID string, r *http.Request) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "hook.handlePlaylist")
	defer span.Finish()

	chunkID := uint64(0)
	if chid, ok := h.segmentsCount.Load(streamID); ok {
		chunkID = chid.(uint64) + 1
		h.segmentsCount.Store(streamID, chunkID)
	}

	path := r.FormValue("path")
	if path == "" {
		return errors.New("failed to get stream path")
	}
	span.SetTag("path", path)

	logger := ctxzap.Extract(ctx).With(zap.String("path", path), zap.String("stream_id", streamID))
	logger.Info("updating playlist")

	h.playlists.Store(streamID, path)

	segments, err := ExtractSegments(path)
	if err != nil {
		return err
	}

	if len(segments) <= 0 {
		return errors.New("no segments")
	}

	stream, err := h.sc.Streams.Get(spanCtx, newStreamRequest(streamID))
	if err != nil {
		return fmt.Errorf("failed to get stream: %s", err)
	}

	if stream.Status == v1.StreamStatusFailed ||
		stream.Status == v1.StreamStatusCancelled ||
		stream.Status == v1.StreamStatusCompleted ||
		stream.Status == v1.StreamStatusDeleted {
		return nil
	}

	achReq := &dispatcherv1.AddInputChunkRequest{
		StreamId:         streamID,
		StreamContractId: stream.StreamContractID,
		ChunkId:          chunkID,
		Reward:           stream.ProfileCost / 60 * segments[chunkID-1].Duration,
	}

	logger.Info("add input chunk", zap.Uint64("chunk_id", chunkID))

	achResp, err := h.sc.Dispatcher.AddInputChunk(spanCtx, achReq)
	if err != nil {
		h.addInputChunkFailed.Store(streamID, chunkID)
		return fmt.Errorf("failed to add input chunk: %s", err)
	}

	logger.Info(
		"add input chunk succesfully",
		zap.String("tx", achResp.Tx),
		zap.String("status", achResp.Status.String()),
		zap.Uint64("chunk_id", chunkID))

	return nil
}

func (h *Hook) handleUpdatePublish(ctx context.Context, streamID string) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "hook.handleUpdatePublish")
	defer span.Finish()

	logger := ctxzap.Extract(ctx).With(zap.String("stream_id", streamID))
	logger.Info("checking publication")

	if i, ok := h.addInputChunkFailed.Load(streamID); ok {
		if i.(uint64) > 0 {
			return fmt.Errorf("failed to add input chunk %d", i.(uint64))
		}
	}

	stream, err := h.sc.Streams.Get(spanCtx, newStreamRequest(streamID))
	if err != nil {
		return fmt.Errorf("failed to get stream: %s", err)
	}

	logger.Info("stream status is", zap.String("status", stream.Status.String()))

	if stream.Status == v1.StreamStatusFailed ||
		stream.Status == v1.StreamStatusCancelled ||
		stream.Status == v1.StreamStatusCompleted {
		return fmt.Errorf("stream status is %s", stream.Status.String())
	}

	return nil
}

func newStreamRequest(id string) *privatev1.StreamRequest {
	return &privatev1.StreamRequest{Id: id}
}
