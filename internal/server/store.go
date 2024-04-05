/*
 * This file was last modified at 2024-04-03 09:42 by Victor N. Skurikhin.
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

func Storage() storage.Storage {
	if store == nil {
		if env.Server.IsDBSetup() {
			store = storage.New(mem, PgxPoolInstance().GetPool())
		} else {
			store = mem
		}
	}
	return store
}
