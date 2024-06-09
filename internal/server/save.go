/*
 * This file was last modified at 2024-06-11 09:57 by Victor N. Skurikhin.
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
