package hookd

import (
	"fmt"
	"time"

	logrussentry "github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Loglevel           string `default:"INFO" envconfig:"LOG_LEVEL"`
	Addr               string `default:"127.0.0.1:8888"`
	UserProfileRPCADDR string `required:"true" default:"127.0.0.1:7001"`
	CamerasRPCADDR     string `required:"true" default:"127.0.0.1:8019"`
	SentryDSN          string `required:"false"`

	Logger *logrus.Entry `ignored:"true"`
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
