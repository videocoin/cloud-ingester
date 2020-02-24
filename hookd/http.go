package hookd

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	privatev1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/hookd/pkg/grpcclient"
	"google.golang.org/grpc"
)

type HTTPServerConfig struct {
	Addr           string
	StreamsRPCAddr string
}

type HTTPServer struct {
	cfg    *HTTPServerConfig
	e      *echo.Echo
	logger *logrus.Entry
	hook   *Hook
}

func NewHTTPServer(
	cfg *HTTPServerConfig,
	logger *logrus.Entry,
	emitter emitterv1.EmitterServiceClient,
) (*HTTPServer, error) {
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
	streams := privatev1.NewStreamsServiceClient(conn)

	hookConfig := &HookConfig{
		Prefix: "/hook",
	}

	hook, err := NewHook(
		e,
		hookConfig,
		streams,
		emitter,
		logger.WithField("system", "hook"),
	)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{
		cfg:    cfg,
		e:      e,
		logger: logger,
		hook:   hook,
	}, nil
}

func (s *HTTPServer) Start() error {
	s.logger.Infof("http server listening on %s", s.cfg.Addr)
	return s.e.Start(s.cfg.Addr)
}

func (s *HTTPServer) Stop() error {
	s.logger.Infof("stopping http server")
	return nil
}
