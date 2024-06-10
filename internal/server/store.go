/*
 * This file was last modified at 2024-06-10 19:19 by Victor N. Skurikhin.
 * store.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/storage"
)

var store storage.Storage

var mem = new(storage.MemStorage)

func init() {
	mem.Metrics = make(map[string]*string)
}

func Storage(cfg env.Config) storage.Storage {
	if store == nil {
		if cfg.IsDBSetup() {
			store = storage.New(mem, PgxPoolInstance().GetPool())
		} else {
			store = mem
		}
	}
	return store
}
