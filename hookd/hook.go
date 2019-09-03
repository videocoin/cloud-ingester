package hookd

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	jobsv1 "github.com/videocoin/cloud-api/jobs/v1"
	managerv1 "github.com/videocoin/cloud-api/manager/v1"
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
	manager managerv1.ManagerServiceClient
}

func NewHook(
	e *echo.Echo,
	cfg *HookConfig,
	manager managerv1.ManagerServiceClient,
	logger *logrus.Entry,
) (*Hook, error) {
	hook := &Hook{
		e:       e,
		cfg:     cfg,
		logger:  logger,
		manager: manager,
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
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "handlePublish")
	defer span.Finish()

	logger := h.logger.WithField("hook", "publish")
	logger.Info("handling hook")

	streamInfo, err := ParseStreamName(spanCtx, r.FormValue("name"))
	if err != nil {
		logger.Warningf("failed to parse stream name: %s", err)
		return ErrBadRequest
	}

	span.SetTag("job_id", streamInfo.JobID)
	logger = logger.WithField("job_id", streamInfo.JobID)

	managerResp, err := h.manager.UpdateStatus(context.Background(), &managerv1.UpdateJobRequest{
		Id:          streamInfo.JobID,
		Status:      jobsv1.JobStatusPending,
		InputStatus: jobsv1.InputStatusActive,
	})

	if err != nil {
		logger.Errorf("failed to update stream status: %s", err.Error())
	}

	logger.Debugf("manager response: %+v", managerResp)

	return nil
}

func (h *Hook) handlePublishDone(ctx context.Context, r *http.Request) error {
	span, spanCtx := opentracing.StartSpanFromContext(ctx, "handlePublishDone")
	defer span.Finish()

	logger := h.logger.WithField("hook", "publish_done")
	logger.Info("handling hook")

	streamInfo, err := ParseStreamName(spanCtx, r.FormValue("name"))
	if err != nil {
		logger.Warningf("failed to parse stream name: %s", err)
		return ErrBadRequest
	}

	span.SetTag("job_id", streamInfo.JobID)
	logger = logger.WithField("job_id", streamInfo.JobID)

	logger.Info("stopping stream")

	managerResp, err := h.manager.Stop(context.Background(), &managerv1.JobRequest{
		Id: streamInfo.JobID,
	})
	if err != nil {
		logger.Errorf("failed to stop stream: %s", err.Error())
	}

	logger.Debugf("manager response: %+v", managerResp)

	return nil
}
