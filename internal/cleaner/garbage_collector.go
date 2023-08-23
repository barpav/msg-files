package cleaner

import (
	"context"
	"fmt"
	"time"

	"github.com/barpav/msg-files/statistics"
)

type GarbageCollector struct {
	cfg     *config
	stats   *statistics.Client
	storage Storage
}

type Storage interface {
	UpdateStats(ctx context.Context, fileId string, inUse bool) (uses int, err error)
	MarkAsUnused(ctx context.Context, fileId string) error
	RemoveFromUnused(ctx context.Context, fileId string) error
	UnusedFiles(ctx context.Context, until, limit int64) (ids []string, err error)
	DeleteFile(ctx context.Context, id string) error
}

func (g *GarbageCollector) Start(storage Storage) (err error) {
	g.cfg = &config{}
	g.cfg.read()

	g.storage = storage

	g.stats = &statistics.Client{}
	err = g.stats.Connect()

	if err == nil {
		err = g.stats.Subscribe(g.UpdateFileUsage)
	}

	if err != nil {
		return fmt.Errorf("failed to start files garbage collector: %w", err)
	}

	go func() {
		for {
			time.Sleep(time.Duration(g.cfg.checkPeriod) * time.Second)
			g.deleteUnusedFiles()
		}
	}()

	return nil
}

func (g *GarbageCollector) Stop(ctx context.Context) (err error) {
	err = g.stats.Disconnect(ctx)

	if err != nil {
		return fmt.Errorf("failed to stop files garbage collector: %w", err)
	}

	return nil
}
