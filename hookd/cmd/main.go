package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gitlab.videocoin.io/videocoin/ingester/hookd"
)

func main() {
	var config hookd.Config

	err := envconfig.Process("LP_STREAMINGESTER_HOOK", &config)
	if err != nil {
		logrus.Fatal(err)
	}

	err = config.InitLogger()
	if err != nil {
		logrus.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		config.Logger.Infof("recieved signal %s", sig)
		done <- true
	}()

	errCh := make(chan error, 1)

	service, err := hookd.NewService(&config)
	if err != nil {
		config.Logger.Fatal(err)
	}

	go func() {
		err := service.Start()
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		config.Logger.Error(err)
	case <-done:
		config.Logger.Info("terminating")
	}

	service.Stop()
}
