package main

import (
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger     *logrus.Entry
	cfg        *Config
	httpServer *httpServer
}

func NewService(cfg *Config) (*Service, error) {
	httpServerCfg := &HTTPServerConfig{
		Addr:               cfg.Addr,
		UserProfileRpcAddr: cfg.UserProfileRpcAddr,
		CamerasRpcAddr:     cfg.CamerasRpcAddr,
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

func (s *Service) Start() error {
	go s.httpServer.Start()
	return nil
}

func (s *Service) Stop() error {
	s.httpServer.Stop()
	return nil
}
