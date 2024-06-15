/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
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
	select {
	case <-ctx.Done():
	default:
		time.Sleep(cfg.StoreInterval() * time.Second)
		store.SaveToFile(cfg.FileStoragePath())
	}
}
