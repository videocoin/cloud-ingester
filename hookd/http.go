package hookd

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	profilesv1 "github.com/videocoin/cloud-api/profiles/manager/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
)

type HTTPServerConfig struct {
	Addr     string
	Emitter  emitterv1.EmitterServiceClient
	Streams  streamsv1.StreamsServiceClient
	Profiles profilesv1.ProfileManagerServiceClient
	Logger   *logrus.Entry
}

type httpServer struct {
	cfg    *HTTPServerConfig
	e      *echo.Echo
	logger *logrus.Entry
	hook   *Hook
}

func NewHTTPServer(
	cfg *HTTPServerConfig,

) (*httpServer, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})
	hookConfig := &HookConfig{
		Prefix: "/hook",
	}

	hook, err := NewHook(
		e,
		hookConfig,
		cfg.Streams,
		cfg.Emitter,
		cfg.Profiles,
		cfg.Logger.WithField("system", "hook"),
	)
	if err != nil {
		return nil, err
	}

	return &httpServer{
		cfg:    cfg,
		e:      e,
		logger: cfg.Logger,
		hook:   hook,
	}, nil
}

func (s *httpServer) Start() error {
	s.logger.Infof("http server listening on %s", s.cfg.Addr)
	return s.e.Start(s.cfg.Addr)
}

func (s *httpServer) Stop() error {
	s.logger.Infof("stopping http server")
	return nil
}
