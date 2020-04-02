package hookd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

type Cleaner struct {
	logger *logrus.Entry
	ticker *time.Ticker
	hlsDir string
}

func NewCleaner(hlsDir string, logger *logrus.Entry) (*Cleaner, error) {
	return &Cleaner{
		logger: logger,
		ticker: time.NewTicker(time.Second * 60),
		hlsDir: hlsDir,
	}, nil
}

func (c *Cleaner) Start() {
	for range c.ticker.C {
		err := c.cleanup()
		if err != nil {
			c.logger.Errorf("failed to cleanup: %s", err)
		}
	}
}

func (c *Cleaner) Stop() error {
	c.ticker.Stop()
	return nil
}

func (c *Cleaner) cleanup() error {
	files, err := ioutil.ReadDir(c.hlsDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if f.ModTime().Before(time.Now().Add(time.Hour * -12)) {
				path := filepath.Join(c.hlsDir, f.Name())
				c.logger.Infof("removing %s", path)
				err := os.RemoveAll(path)
				if err != nil {
					c.logger.Errorf("failed to remove %s", path)
					continue
				}
			}
		}
	}

	return nil
}
