package hookd

import (
	"bufio"
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/grafov/m3u8"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	v1 "github.com/videocoin/cloud-api/streams/v1"
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
	e             *echo.Echo
	logger        *logrus.Entry
	streams       privatev1.StreamsServiceClient
	emitter       emitterv1.EmitterServiceClient
	segmentsCount sync.Map
}

func NewHook(
	e *echo.Echo,
	cfg *HookConfig,
	streams privatev1.StreamsServiceClient,
	emitter emitterv1.EmitterServiceClient,
	logger *logrus.Entry,
) (*Hook, error) {
	hook := &Hook{
		e:       e,
		cfg:     cfg,
		logger:  logger,
		streams: streams,
		emitter: emitter,
	}
	hook.e.Any(cfg.Prefix, hook.handleHook)
	return hook, nil
}

func (h *Hook) handleHook(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()

	err := req.ParseForm()
	if err != nil {
		h.logger.Errorf("failed to parse form: %s", err)
		return ErrBadRequest
	}

	call := req.FormValue("call")

	h.logger.Infof("hook %s", call)
	h.logger.Infof("hook params %+v", req.Form)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "handlePublish")
	defer span.Finish()

	logger := h.logger.WithField("hook", "publish")
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warningf("failed to get stream id")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	logger = logger.WithField("id", streamID)

	req := &privatev1.StreamRequest{Id: streamID}
	streamResp, err := h.streams.Get(ctx, req)
	if err != nil {
		logger.Errorf("failed to get stream: %s", err.Error())
		return ErrBadRequest
	}

	if streamResp.Status != v1.StreamStatusPrepared {
		logger.Errorf("wrong stream status: %s", streamResp.Status.String())
		return ErrBadRequest
	}

	streamResp, err = h.streams.Publish(ctx, req)
	if err != nil {
		logger.Errorf("failed to publish: %s", err.Error())
		return ErrBadRequest
	}

	return nil
}

func (h *Hook) handlePublishDone(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handlePublishDone")
	defer span.Finish()

	logger := h.logger.WithField("hook", "publish_done")
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warningf("failed to get stream id")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	logger = logger.WithField("id", streamID)

	logger.Info("publishing done")

	req := &privatev1.StreamRequest{Id: streamID}
	_, err := h.streams.PublishDone(ctx, req)
	if err != nil {
		logger.Errorf("failed to publish done: %s", err.Error())
		return ErrBadRequest
	}

	return nil
}

func (h *Hook) handlePlaylist(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handlePlaylist")
	defer span.Finish()

	logger := h.logger.WithField("hook", "playlist")
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warningf("failed to get stream name")
		return ErrBadRequest
	}

	path := r.FormValue("path")
	if path == "" {
		logger.Warningf("failed to get stream path")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	span.SetTag("path", path)
	logger = logger.WithField("id", streamID)

	logger.Info("updating playlist")

	f, err := os.Open(path)
	if err != nil {
		logger.Errorf("failed to open playlist: %s", err)
		return ErrInternalServer
	}

	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		logger.Errorf("failed to decode playlist: %s", err)
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
		streamResp, err := h.streams.Get(ctx, req)
		if err != nil {
			logger.Errorf("failed to get stream: %s", err.Error())
			return ErrBadRequest
		}

		actual, ok := h.segmentsCount.LoadOrStore(streamID, segmentsCount)
		if ok {
			h.segmentsCount.Store(streamID, segmentsCount)
		}
		prevSegmentsCount := actual.(int)

		for i := prevSegmentsCount; i < segmentsCount; i++ {
			achReq := &emitterv1.AddInputChunkIdRequest{
				StreamContractId: streamResp.StreamContractID,
				ChunkId:          uint64(i),
				ChunkDuration:    pl.Segments[i-1].Duration,
			}

			logger.WithFields(logrus.Fields{
				"stream_contract_id": streamResp.StreamContractID,
				"chunk_id":           i,
			}).Debugf("calling AddInputChunkId")

			_, err := h.emitter.AddInputChunkId(ctx, achReq)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"stream_contract_id": streamResp.StreamContractID,
					"chunk_id":           i,
				}).Errorf("failed to add input chunk: %s", err.Error())
			}
		}

		if pl.Closed {
			_, err := h.streams.PublishDone(ctx, req)
			if err != nil {
				logger.Errorf("failed to publish done: %s", err.Error())
				return nil
			}
		}
	}

	return nil
}

func (h *Hook) handleUpdatePublish(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handleUpdatePublish")
	defer span.Finish()

	logger := h.logger.WithField("hook", "update")
	logger.Info("handling hook")

	streamID := r.FormValue("name")
	if streamID == "" {
		logger.Warningf("failed to get stream id")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamID)
	logger = logger.WithField("id", streamID)

	req := &privatev1.StreamRequest{Id: streamID}
	streamResp, err := h.streams.Get(ctx, req)
	if err != nil {
		logger.Errorf("failed to get stream: %s", err.Error())
		return nil
	}

	logger.Infof("stream status is %s", streamResp.Status.String())

	if streamResp.Status == v1.StreamStatusFailed ||
		streamResp.Status == v1.StreamStatusCancelled ||
		streamResp.Status == v1.StreamStatusCompleted {
		return ErrBadRequest
	}

	return nil
}
