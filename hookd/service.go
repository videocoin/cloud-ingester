package hookd

import (
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"google.golang.org/grpc"
)

type Service struct {
	logger     *logrus.Entry
	cfg        *Config
	httpServer *httpServer
	cleaner    *Cleaner
}

func NewService(cfg *Config) (*Service, error) {
	elogger := cfg.Logger.WithField("system", "emittercli")
	eGrpcDialOpts := grpcutil.ClientDialOptsWithRetry(elogger)
	emitterConn, err := grpc.Dial(cfg.EmitterRPCAddr, eGrpcDialOpts...)
	if err != nil {
		return nil, err
	}
	emitter := emitterv1.NewEmitterServiceClient(emitterConn)

	httpServerCfg := &HTTPServerConfig{
		Addr:           cfg.Addr,
		StreamsRPCAddr: cfg.StreamsRPCAddr,
	}
	httpServer, err := NewHTTPServer(
		httpServerCfg,
		cfg.Logger.WithField("system", "http-server"),
		emitter,
	)
	if err != nil {
		return nil, err
	}

	cleaner, err := NewCleaner(cfg.HLSDir, cfg.Logger.WithField("system", "cleaner"))
	if err != nil {
		return nil, err
	}

	return &Service{
		logger:     cfg.Logger,
		cfg:        cfg,
		httpServer: httpServer,
		cleaner:    cleaner,
	}, nil
}

func (s *Service) Start() error {
	go s.httpServer.Start()
	go s.cleaner.Start()
	return nil
}

func (s *Service) Stop() error {
	s.cleaner.Stop()
	s.httpServer.Stop()

	return nil
}
