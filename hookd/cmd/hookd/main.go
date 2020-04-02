package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-ingester/hookd"
)

func main() {
	var config hookd.Config

	err := envconfig.Process("VC_HOOKD", &config)
	if err != nil {
		logrus.Fatal(err)
	}

	err = config.InitLogger()
	if err != nil {
		logrus.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)
	errCh := make(chan error, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		config.Logger.Infof("recieved signal %s", sig)
		exit <- true
	}()

	service, err := hookd.NewService(&config)
	if err != nil {
		config.Logger.Fatal(err)
	}

	go service.Start(errCh)

	select {
	case <-exit:
		break
	case err := <-errCh:
		if err != nil {
			config.Logger.Error(err)
		}
		break
	}

	err = service.Stop()
	if err != nil {
		logrus.Fatal(err)
	}
}
