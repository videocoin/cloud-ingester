package hookd

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	v1 "github.com/videocoin/cloud-api/streams/v1"
	"github.com/videocoin/hookd/pkg/grpcclient"
	"google.golang.org/grpc"
)

type HTTPServerConfig struct {
	Addr           string
	StreamsRPCAddr string
}

type httpServer struct {
	cfg    *HTTPServerConfig
	e      *echo.Echo
	logger *logrus.Entry
	hook   *Hook
}

func NewHTTPServer(cfg *HTTPServerConfig, logger *logrus.Entry) (*httpServer, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	opts := grpcclient.DialOpts(logger)
	conn, err := grpc.Dial(cfg.StreamsRPCAddr, opts...)
	if err != nil {
		return nil, err
	}
	streams := v1.NewStreamServiceClient(conn)

	hookConfig := &HookConfig{
		Prefix: "/hook",
	}

	hook, err := NewHook(
		e,
		hookConfig,
		streams,
		logger.WithField("system", "hook"),
	)
	if err != nil {
		return nil, err
	}

	return &httpServer{
		cfg:    cfg,
		e:      e,
		logger: logger,
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
