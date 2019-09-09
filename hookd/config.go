package hookd

import (
	"fmt"
	"time"

	logrussentry "github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Addr           string `required:"true" envconfig:"ADDR" default:"0.0.0.0:8887"`
	StreamsRPCAddr string `required:"true" envconfig:"STREAMS_RPC_ADDR" default:"127.0.0.1:50051"`

	SentryDSN string        `required:"false"`
	Logger    *logrus.Entry `ignored:"true"`
	Loglevel  string        `required:"true" envconfig:"LOGLEVEL" default:"INFO"`
}

func (c *Config) InitLogger() error {
	level, err := logrus.ParseLevel(c.Loglevel)
	if err != nil {
		return fmt.Errorf("not a valid log level: %q", c.Loglevel)
	}

	logrus.SetLevel(level)

	if level == logrus.DebugLevel {
		logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	}

	if c.SentryDSN != "" {
		sentryLevels := []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		}
		sentryTags := map[string]string{
			"service": Name,
			"version": Version,
		}
		sentryHook, err := logrussentry.NewWithTagsSentryHook(
			c.SentryDSN,
			sentryTags,
			sentryLevels,
		)
		sentryHook.StacktraceConfiguration.Enable = true
		sentryHook.Timeout = 5 * time.Second
		sentryHook.SetRelease(Version)

		if err == nil {
			logrus.AddHook(sentryHook)
		} else {
			logrus.Warning(err)
		}
	}

	logger := logrus.WithFields(logrus.Fields{
		"service": Name,
		"version": Version,
	})

	c.Logger = logger

	return nil
}
