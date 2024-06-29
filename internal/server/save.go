/*
 * This file was last modified at 2024-06-25 00:02 by Victor N. Skurikhin.
 * save.go
 * $Id$
 */

package server

import (
	"context"
	"github.com/vskurikhin/gometrics/internal/env"
	"time"
)

func SaveLoop(ctx context.Context, cfg env.Config) {
	select {
	case <-ctx.Done():
	default:
		time.Sleep(cfg.StoreInterval())
		Save(cfg)
	}
}

func Save(cfg env.Config) {
	if cfg.Restore() && cfg.FileStoragePath() != "" {
		store.SaveToFile(cfg.FileStoragePath())
	}
}
