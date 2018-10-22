package hookd

import (
	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"gitlab.videocoin.io/ingester/hookd/pkg/grpcclient"
	pb "gitlab.videocoin.io/ingester/hookd/pkg/liveplanet/api/proto"
	"google.golang.org/grpc"
)

// HTTPServerConfig addresses for http server
type HTTPServerConfig struct {
	Addr               string
	UserProfileRPCADDR string
	CamerasRPCADDR     string
}

type httpServer struct {
	cfg    *HTTPServerConfig
	e      *echo.Echo
	logger *logrus.Entry
	hook   *Hook
}

// NewHTTPServer returns reference to new httpServer object
func NewHTTPServer(cfg *HTTPServerConfig, logger *logrus.Entry) (*httpServer, error) {
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

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/metrics", echo.WrapHandler(prometheus.Handler()))
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	hook, err := NewHook(
		e,
		"/hook",
		profiles,
		cameras,
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
