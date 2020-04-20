package cleaner

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

type Cleaner struct {
	logger *zap.Logger
	ticker *time.Ticker
	hlsDir string
}

func NewCleaner(ctx context.Context, hlsDir string) (*Cleaner, error) {
	return &Cleaner{
		logger: ctxzap.Extract(ctx).With(zap.String("system", "cleaner")),
		ticker: time.NewTicker(time.Second * 60),
		hlsDir: hlsDir,
	}, nil
}

func (c *Cleaner) Start() {
	for range c.ticker.C {
		err := c.cleanup()
		if err != nil {
			c.logger.Error("failed to cleanup", zap.Error(err))
		}
	}
}

func (c *Cleaner) Stop() {
	c.ticker.Stop()
}

func (c *Cleaner) cleanup() error {
	if _, err := os.Stat(c.hlsDir); os.IsNotExist(err) {
		return nil
	}

	files, err := ioutil.ReadDir(c.hlsDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			if f.ModTime().Before(time.Now().Add(time.Hour * -12)) {
				path := filepath.Join(c.hlsDir, f.Name())
				c.logger.Info("removing", zap.String("path", path))
				err := os.RemoveAll(path)
				if err != nil {
					c.logger.Error("failed to remove", zap.String("path", path))
					continue
				}
			}
		}
	}

	return nil
}
