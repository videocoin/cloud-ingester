package hookd

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	managerv1 "github.com/videocoin/cloud-api/manager/v1"
	"github.com/videocoin/cloud-ingester/hookd/pkg/grpcclient"
	"google.golang.org/grpc"
)

type HTTPServerConfig struct {
	Addr           string
	ManagerRPCAddr string
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

	hookConfig := &HookConfig{Prefix: "/hook"}

	managerOpts := grpcclient.DialOpts(logger)
	managerConn, err := grpc.Dial(cfg.ManagerRPCAddr, managerOpts...)
	if err != nil {
		return nil, err
	}
	manager := managerv1.NewManagerServiceClient(managerConn)

	hook, err := NewHook(
		e,
		hookConfig,
		manager,
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
