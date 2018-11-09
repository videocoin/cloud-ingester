package hookd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Service struct used for hookd service object
type Service struct {
	logger     *log.Entry
	cfg        *Config
	HTTPServer *HTTPServer
}

// NewService returns new	ingest hook service
func NewService(cfg *Config) (*Service, error) {
	httpServerCfg := &HTTPServerConfig{
		Addr:           cfg.Addr,
		ManagerRPCADDR: cfg.ManagerRPCADDR,
	}
	HTTPServer, err := NewHTTPServer(
		httpServerCfg,
		log.WithField("system", "http-server"),
	)
	if err != nil {
		return nil, err
	}
	return &Service{
		logger:     log.New().WithField("name", "hookd"),
		cfg:        cfg,
		HTTPServer: HTTPServer,
	}, nil
}

// Start runs hookd service
func Start() {
	cfgLoc, exists := os.LookupEnv("CONFIG_LOC")
	if !exists {
		log.Fatalf("please set CONFIG_LOC")
	}
	cfg := LoadConfig(cfgLoc)

	level, err := log.ParseLevel(cfg.Loglevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(level)

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		log.Infof("recieved signal %s", sig)
		done <- true
	}()

	errCh := make(chan error, 1)

	service, err := NewService(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := service.StartHTTP()
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		log.Error(err)
	case <-done:
		log.Info("terminating")
	}

	service.StopHTTP()

}

// Start runs http server
func (s *Service) StartHTTP() error {
	go s.HTTPServer.Start()
	return nil
}

// Stop stops the http server from service
func (s *Service) StopHTTP() error {
	s.HTTPServer.Stop()
	return nil
}
