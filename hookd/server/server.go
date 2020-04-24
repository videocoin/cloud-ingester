package server

import (
	"context"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/labstack/echo/v4"
	otecho "github.com/opentracing-contrib/echo"
	echologrus "github.com/plutov/echo-logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
)

type Server struct {
	logger *logrus.Entry
	addr   string
	e      *echo.Echo
	hook   *Hook
}

func NewServer(ctx context.Context, addr string, sc *clientv1.ServiceClient) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	sn, _ := ServiceNameFromContext(ctx)

	echologrus.Logger = ctxlogrus.Extract(ctx).Logger
	e.Logger = echologrus.GetEchoLogger()

	e.Use(otecho.Middleware(sn))
	e.Use(echologrus.Hook())

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	hookConfig := &HookConfig{
		Prefix: "/hook",
	}

	logger := grpclogrus.Extract(ctx).WithField("system", "server")

	hook, err := NewHook(ctxlogrus.ToContext(ctx, logger), e, hookConfig, sc)
	if err != nil {
		return nil, err
	}

	return &Server{
		logger: logger,
		addr:   addr,
		e:      e,
		hook:   hook,
	}, nil
}

func (s *Server) Start() error {
	return s.e.Start(s.addr)
}

func (s *Server) Stop() error {
	return s.e.Shutdown(context.Background())
}
