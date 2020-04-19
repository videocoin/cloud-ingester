package server

import (
	"context"

	"github.com/brpaz/echozap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/labstack/echo/v4"
	otecho "github.com/opentracing-contrib/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
	addr   string
	e      *echo.Echo
	hook   *Hook
}

func NewServer(ctx context.Context, addr string, sc *clientv1.ServiceClient) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	sn, _ := ServiceNameFromContext(ctx)

	e.Use(otecho.Middleware(sn))
	e.Use(echozap.ZapLogger(ctxzap.Extract(ctx)))

	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	hookConfig := &HookConfig{
		Prefix: "/hook",
	}

	hook, err := NewHook(ctx, e, hookConfig, sc)
	if err != nil {
		return nil, err
	}

	return &Server{
		logger: ctxzap.Extract(ctx).With(zap.String("system", "server")),
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
