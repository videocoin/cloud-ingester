package hookd

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	pb "gitlab.videocoin.io/videocoin/common/proto"
	"gitlab.videocoin.io/videocoin/ingester/hookd/pkg/grpcclient"
	"google.golang.org/grpc"
)

// HTTPServerConfig addresses for http server
type HTTPServerConfig struct {
	Addr               string
	UserProfileRPCADDR string
	CamerasRPCADDR     string
	ManagerRPCADDR     string
}

// HTTPServer http server reciver
// holds echo, config, and logger objects
type HTTPServer struct {
	cfg    *HTTPServerConfig
	e      *echo.Echo
	logger *logrus.Entry
	hook   *Hook
}

// NewHTTPServer returns reference to new HTTPServer object
func NewHTTPServer(cfg *HTTPServerConfig, logger *logrus.Entry) (*HTTPServer, error) {
	grpcDialOpts := grpcclient.DialOpts(logger)
	upConn, err := grpc.Dial(cfg.UserProfileRPCADDR, grpcDialOpts...)
	if err != nil {
		return nil, err
	}
	profiles := pb.NewUserProfileServiceClient(upConn)

	camerasConn, err := grpc.Dial(cfg.CamerasRPCADDR, grpcDialOpts...)
	if err != nil {
		return nil, err
	}
	cameras := pb.NewCameraCloudInternalServiceClient(camerasConn)

	managerConn, err := grpc.Dial(cfg.ManagerRPCADDR, grpcDialOpts...)
	if err != nil {
		return nil, err
	}

	manager := pb.NewManagerServiceClient(managerConn)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/metrics", echo.WrapHandler(prometheus.Handler()))
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	status, err := manager.Health(context.Background(), &empty.Empty{})
	if status.String() != "healthy" || err != nil {
		panic(err)
	}

	hook, err := NewHook(
		e,
		"/hook",
		profiles,
		manager,
		cameras,
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

// Start starts echo server
func (s *HTTPServer) Start() error {
	s.logger.Infof("http server listening on %s", s.cfg.Addr)
	return s.e.Start(s.cfg.Addr)
}

// Stop does nothing
func (s *HTTPServer) Stop() error {
	s.logger.Infof("stopping http server")
	return nil
}
