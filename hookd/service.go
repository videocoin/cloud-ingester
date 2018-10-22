package hookd

import (
	"github.com/sirupsen/logrus"
)

// Service struct used for hookd service object
type Service struct {
	logger     *logrus.Entry
	cfg        *Config
	httpServer *httpServer
}

// NewService returns new	ingest hook service
func NewService(cfg *Config) (*Service, error) {
	httpServerCfg := &HTTPServerConfig{
		Addr:               cfg.Addr,
		UserProfileRPCADDR: cfg.UserProfileRPCADDR,
		CamerasRPCADDR:     cfg.CamerasRPCADDR,
	}
	httpServer, err := NewHTTPServer(
		httpServerCfg,
		cfg.Logger.WithField("system", "http-server"),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		logger:     cfg.Logger,
		cfg:        cfg,
		httpServer: httpServer,
	}, nil
}

// Start runs http server
func (s *Service) Start() error {
	go s.httpServer.Start()
	return nil
}

// Stop stops the http server from service
func (s *Service) Stop() error {
	s.httpServer.Stop()
	return nil
}
