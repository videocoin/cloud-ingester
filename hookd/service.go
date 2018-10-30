package hookd

import (
	"github.com/sirupsen/logrus"
)

// Service struct used for hookd service object
type Service struct {
	logger     *logrus.Entry
	cfg        *Config
	HTTPServer *HTTPServer
}

// NewService returns new	ingest hook service
func NewService(cfg *Config) (*Service, error) {
	httpServerCfg := &HTTPServerConfig{
		Addr:               cfg.Addr,
		UserProfileRPCADDR: cfg.UserProfileRPCADDR,
		CamerasRPCADDR:     cfg.CamerasRPCADDR,
		ManagerRPCADDR:     cfg.ManagerRPCADDR,
	}
	HTTPServer, err := NewHTTPServer(
		httpServerCfg,
		cfg.Logger.WithField("system", "http-server"),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		logger:     cfg.Logger,
		cfg:        cfg,
		HTTPServer: HTTPServer,
	}, nil
}

// Start runs http server
func (s *Service) Start() error {
	go s.HTTPServer.Start()
	return nil
}

// Stop stops the http server from service
func (s *Service) Stop() error {
	s.HTTPServer.Stop()
	return nil
}
