/*
 * This file was last modified at 2024-03-18 22:49 by Victor N. Skurikhin.
 * store.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"github.com/vskurikhin/gometrics/internal/storage/postgres"
)

var store storage.Storage

var mem = new(memory.MemStorage)

func init() {
	mem.Metrics = make(map[string]*string)
}

func Storage() storage.Storage {
	if store == nil {
		if env.Server.IsDBSetup() {
			store = postgres.New(mem, PgxPoolInstance().GetPool())
		} else {
			store = mem
		}
	}
	return store
}
