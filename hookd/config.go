package hookd

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config all required config for hookd project
type Config struct {
	Addr           string `required:"true" default:"127.0.0.1:8888"`
	Loglevel       string `required:"true" default:"DEBUG" envconfig:"LOG_LEVEL"`
	ManagerRPCADDR string `required:"true" default:"127.0.0.1:50051"`
	SentryDSN      string `required:"false"`
}

var cfg Config
var once sync.Once

// LoadConfig initialize config
func LoadConfig(loc string) *Config {
	switch loc {
	case "local":
		once.Do(func() {
			err := envconfig.Process("INGESTER", &cfg)
			if err != nil {
				logrus.Fatalf("failed to load config: %s", err.Error())
			}
		})
		break
	// requires PROJECT_ID environment variable
	case "remote":
		once.Do(func() {
			ctx := context.Background()
			client, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
			if err != nil {
				logrus.Fatalf("failed to create new client: %s", err)
			}

			key := datastore.NameKey("config", "ingester", nil)
			err = client.Get(ctx, key, &cfg)
			if err != nil {
				logrus.Fatalf("failed to get namekey: %s", err)
			}
		})

		break

	default:

	}

	return &cfg
}
