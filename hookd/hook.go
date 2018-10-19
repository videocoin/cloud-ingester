package hookd

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	pb "gitlab.videocoin.io/ingester/hookd/pkg/liveplanet/api/proto"
)

var (
	ErrUnknownHook = errors.New("unknown hook")

	ErrBadRequest = echo.NewHTTPError(http.StatusBadRequest)
)

type Hook struct {
	e       *echo.Echo
	logger  *logrus.Entry
	profile pb.UserProfileServiceClient
	cameras pb.CameraCloudInternalServiceClient
}

func NewHook(
	e *echo.Echo,
	prefix string,
	profile pb.UserProfileServiceClient,
	cameras pb.CameraCloudInternalServiceClient,
	logger *logrus.Entry,
) (*Hook, error) {
	hook := &Hook{
		e:       e,
		profile: profile,
		cameras: cameras,
		logger:  logger,
	}
	hook.e.Any(prefix, hook.handleHook)
	return hook, nil
}

func (h *Hook) handleHook(c echo.Context) error {
	req := c.Request()

	err := req.ParseForm()
	if err != nil {
		h.logger.Error("failed to parse form: %s", err)
		return ErrBadRequest
	}

	h.logger.Debugf("handle hook %+v", req.Form)

	call := req.FormValue("call")
	switch call {
	case "publish":
		err = h.handlePublish(req)
	case "update_publish":
		err = h.handleUpdatePublish(req)
	case "publish_done":
		err = h.handlePublishDone(req)
	case "record":
		err = h.handleRecord(req)
	case "record_done":
		err = h.handleRecordDone(req)
	default:
		return c.NoContent(http.StatusBadRequest)
	}

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Hook) handlePublish(r *http.Request) error {
	logger := h.logger.WithField("hook", "publish")
	logger.Info("handling hook")

	streamInfo, err := ParseStreamName(r.FormValue("name"))
	if err != nil {
		logger.Warningf("failed to parse stream name: %s", err)
		return ErrBadRequest
	}

	logger = logger.WithFields(logrus.Fields{
		"uid": streamInfo.UserID,
		"cid": streamInfo.CameraID,
	})

	logger.Info("getting user profile")

	ctx := context.Background()
	tokenReq := &pb.OAuth2TokenRequest{
		UserId: streamInfo.UserID,
		AppId:  "web",
	}
	tokenResp, err := h.profile.GetOAuth2Token(ctx, tokenReq)
	if err != nil {
		logger.Errorf("failed to get oath2 token: %s", err)
		return ErrBadRequest
	}

	logger.Debugf("token response: %+v", tokenResp)

	logger.Info("marking camera as on air")

	cameraReq := &pb.InternalCameraRequest{
		ID:      streamInfo.CameraID,
		OwnerID: tokenResp.UserId,
	}
	cameraResp, err := h.cameras.MarkCameraAsOnAir(ctx, cameraReq)
	if err != nil {
		logger.Errorf("failed to mark camera as on air: %s", err)
		return ErrBadRequest
	}

	logger.Debugf("camera response: %+v", cameraResp)

	return nil
}

func (h *Hook) handleUpdatePublish(r *http.Request) error {
	logger := h.logger.WithField("hook", "update_publish")
	logger.Info("handling hook")

	streamInfo, err := ParseStreamName(r.FormValue("name"))
	if err != nil {
		logger.Warningf("failed to parse stream name: %s", err)
		return ErrBadRequest
	}

	logger = logger.WithFields(logrus.Fields{
		"uid": streamInfo.UserID,
		"cid": streamInfo.CameraID,
	})

	logger.Info("getting user profile")

	ctx := context.Background()
	tokenReq := &pb.OAuth2TokenRequest{
		UserId: streamInfo.UserID,
		AppId:  "web",
	}
	tokenResp, err := h.profile.GetOAuth2Token(ctx, tokenReq)
	if err != nil {
		logger.Errorf("failed to get oath2 token: %s", err)
		return ErrBadRequest
	}

	logger.Debugf("token response: %+v", tokenResp)

	logger.Info("getting camera")

	cameraReq := &pb.InternalCameraRequest{
		ID:      streamInfo.CameraID,
		OwnerID: tokenResp.UserId,
	}
	cameraResp, err := h.cameras.GetCamera(ctx, cameraReq)
	if err != nil {
		logger.Errorf("failed to get camera: %s", err)
		return ErrBadRequest
	}

	logger.Debugf("camera response: %+v", cameraResp)

	return nil
}

func (h *Hook) handlePublishDone(r *http.Request) error {
	logger := h.logger.WithField("hook", "publish_done")
	logger.Info("handling hook")

	streamInfo, err := ParseStreamName(r.FormValue("name"))
	if err != nil {
		logger.Warningf("failed to parse stream name: %s", err)
		return ErrBadRequest
	}

	logger = logger.WithFields(logrus.Fields{
		"uid": streamInfo.UserID,
		"cid": streamInfo.CameraID,
	})

	logger.Info("getting user profile")

	ctx := context.Background()
	tokenReq := &pb.OAuth2TokenRequest{
		UserId: streamInfo.UserID,
		AppId:  "web",
	}
	tokenResp, err := h.profile.GetOAuth2Token(ctx, tokenReq)
	if err != nil {
		logger.Errorf("failed to get oath2 token: %s", err)
		return ErrBadRequest
	}

	logger.Debugf("token response: %+v", tokenResp)

	logger.Info("marking camera as off air")

	cameraReq := &pb.InternalCameraRequest{
		ID:      streamInfo.CameraID,
		OwnerID: tokenResp.UserId,
	}
	cameraResp, err := h.cameras.MarkCameraAsOffAir(ctx, cameraReq)
	if err != nil {
		logger.Errorf("failed to mark camera as off air: %s", err)
		return ErrBadRequest
	}

	logger.Debugf("camera response: %+v", cameraResp)

	return nil
}

func (h *Hook) handleRecord(r *http.Request) error {
	return nil
}

func (h *Hook) handleRecordDone(r *http.Request) error {
	return nil
}
