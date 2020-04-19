package service

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	clientv1 "github.com/videocoin/cloud-api/client/v1"
	"github.com/videocoin/cloud-ingester/hookd/cleaner"
	"github.com/videocoin/cloud-ingester/hookd/server"
	"go.uber.org/zap"
)

type Service struct {
	cfg     *Config
	logger  *zap.Logger
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
		logger:  ctxzap.Extract(ctx),
		srv:     srv,
		cleaner: cleaner,
	}, nil
}

func (s *Service) Start(errCh chan error) {
	go func() {
		s.logger.Info("starting server", zap.String("addr", s.cfg.Addr))
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
