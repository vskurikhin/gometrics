/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * save.go
 * $Id$
 */

package server

import (
	"context"
	"time"

	"github.com/vskurikhin/gometrics/internal/env"
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
		cfg.
			Property().
			Storage().
			SaveToFile(cfg.FileStoragePath())
	}
}
