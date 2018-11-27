package hookd

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// Config all required config for hookd project
type Config struct {
	Addr           string `required:"true" default:"ingester:8888"`
	Loglevel       string `required:"true" default:"DEBUG" envconfig:"LOG_LEVEL"`
	ManagerRPCADDR string `required:"true" envconfig:"MANAGER_RPC_ADDR" default:"manager:50051"`
	SentryDSN      string `required:"false"`
}

var cfg Config
var once sync.Once

// LoadConfig initialize config
func LoadConfig() *Config {
	once.Do(func() {
		err := envconfig.Process("", &cfg)
		if err != nil {
			logrus.Fatalf("failed to load config: %s", err.Error())
		}
	})

	return &cfg
}
