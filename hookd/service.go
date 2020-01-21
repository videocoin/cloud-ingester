package hookd

import (
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	profilesv1 "github.com/videocoin/cloud-api/profiles/manager/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
)

type Service struct {
	logger     *logrus.Entry
	cfg        *Config
	httpServer *httpServer
	cleaner    *Cleaner
}

func NewService(cfg *Config) (*Service, error) {
	conn, err := grpcutil.Connect(cfg.StreamsRPCAddr, cfg.Logger.WithField("system", "streamscli"))
	if err != nil {
		return nil, err
	}
	streams := streamsv1.NewStreamsServiceClient(conn)

	conn, err = grpcutil.Connect(cfg.EmitterRPCAddr, cfg.Logger.WithField("system", "emittercli"))
	if err != nil {
		return nil, err
	}
	emitter := emitterv1.NewEmitterServiceClient(conn)

	conn, err = grpcutil.Connect(cfg.ProfilesManagerRPCAddr, cfg.Logger.WithField("system", "profilescli"))
	if err != nil {
		return nil, err
	}
	profiles := profilesv1.NewProfileManagerServiceClient(conn)

	httpServerCfg := &HTTPServerConfig{
		Addr:     cfg.Addr,
		Streams:  streams,
		Emitter:  emitter,
		Profiles: profiles,
		Logger:   cfg.Logger.WithField("system", "http-server"),
	}
	httpServer, err := NewHTTPServer(httpServerCfg)
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
