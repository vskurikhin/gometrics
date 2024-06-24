/*
 * This file was last modified at 2024-06-24 22:51 by Victor N. Skurikhin.
 * save.go
 * $Id$
 */

package server

import (
	"context"
	"github.com/vskurikhin/gometrics/internal/env"
	"time"
)

func Save(ctx context.Context, cfg env.Config) {
	for {
		select {
		case <-ctx.Done():
		default:
			time.Sleep(cfg.StoreInterval())
			if cfg.Restore() && cfg.FileStoragePath() != "" {
				store.SaveToFile(cfg.FileStoragePath())
			}
		}
	}
}
