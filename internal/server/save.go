/*
 * This file was last modified at 2024-06-10 17:37 by Victor N. Skurikhin.
 * save.go
 * $Id$
 */

package server

import (
	"time"

	"github.com/vskurikhin/gometrics/internal/env"
)

func Save(cfg env.Config) {
	for {
		save(cfg)
	}
}

func save(cfg env.Config) {
	time.Sleep(cfg.StoreInterval() * time.Second)
	store.SaveToFile(cfg.FileStoragePath())
}
