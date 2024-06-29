/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * store.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/storage"
	"sync"
)

var (
	mem   = new(storage.MemStorage)
	once  = new(sync.Once)
	store storage.Storage
)

func Storage(cfg env.Config) storage.Storage {
	once.Do(func() {
		mem.Metrics = make(map[string]*string)
		if cfg.IsDBSetup() {
			store = storage.New(mem, pgxPoolInstance().getPool())
		} else {
			store = mem
		}
	})
	return store
}
