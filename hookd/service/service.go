package service

import (
	"context"

	grpclogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-ingester/hookd/cleaner"
	"github.com/videocoin/cloud-ingester/hookd/server"
)

type Service struct {
	cfg     *Config
	logger  *logrus.Entry
	srv     *server.Server
	cleaner *cleaner.Cleaner
}

func NewService(ctx context.Context, cfg *Config) (*Service, error) {
	sc, err := clientv1.NewServiceClientFromEnvconfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	ctx = server.NewContextWithServiceName(ctx, cfg.Name)
	srv, err := server.NewServer(ctx, cfg.Addr, sc)
	if err != nil {
		return nil, err
	}

	cleaner, err := cleaner.NewCleaner(ctx, cfg.HLSDir)
	if err != nil {
		return nil, err
	}

	return &Service{
		cfg:     cfg,
		logger:  grpclogrus.Extract(ctx),
		srv:     srv,
		cleaner: cleaner,
	}, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		s.logger.WithField("addr", s.cfg.Addr).Info("starting server")
		errCh <- s.srv.Start()
	}()
	go func() {
		s.logger.Info("starting cleaner")
		s.cleaner.Start()
	}()
}

func (s *Service) Stop() error {
	s.logger.Info("stopping cleaner")
	s.cleaner.Stop()

	s.logger.Info("stopping server")
	err := s.srv.Stop()
	if err != nil {
		return err
	}

	return nil
}
