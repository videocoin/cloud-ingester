package hookd

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/streams/v1"
)

var (
	ErrBadRequest = echo.NewHTTPError(http.StatusBadRequest)
)

type HookConfig struct {
	Prefix string
}

type Hook struct {
	cfg     *HookConfig
	e       *echo.Echo
	logger  *logrus.Entry
	streams v1.StreamServiceClient
}

func NewHook(
	e *echo.Echo,
	cfg *HookConfig,
	streams v1.StreamServiceClient,
	logger *logrus.Entry,
) (*Hook, error) {
	hook := &Hook{
		e:       e,
		cfg:     cfg,
		logger:  logger,
		streams: streams,
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

	h.logger.Debugf("handle hook %+v", req.Form)

	call := req.FormValue("call")
	switch call {
	case "publish":
		err = h.handlePublish(ctx, req)
	case "publish_done":
		err = h.handlePublishDone(ctx, req)
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

	streamId := r.FormValue("name")
	if streamId == "" {
		logger.Warningf("failed to get stream name")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamId)
	logger = logger.WithField("id", streamId)

	_, err := h.streams.Update(ctx, &v1.UpdateStreamRequest{
		Id:          streamId,
		Status:      v1.StreamStatusPending,
		InputStatus: v1.InputStatusActive,
	})

	if err != nil {
		logger.Errorf("failed to update stream status: %s", err.Error())
	}

	return nil
}

func (h *Hook) handlePublishDone(ctx context.Context, r *http.Request) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "handlePublishDone")
	defer span.Finish()

	logger := h.logger.WithField("hook", "publish_done")
	logger.Info("handling hook")

	streamId := r.FormValue("name")
	if streamId == "" {
		logger.Warningf("failed to get stream name")
		return ErrBadRequest
	}

	span.SetTag("stream_id", streamId)
	logger = logger.WithField("id", streamId)

	logger.Info("stopping stream")

	_, err := h.streams.Stop(ctx, &v1.StreamRequest{
		Id: streamId,
	})

	if err != nil {
		logger.Errorf("failed to stop stream: %s", err.Error())
	}

	return nil
}
