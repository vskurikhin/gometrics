/*
 * This file was last modified at 2024-03-18 23:07 by Victor N. Skurikhin.
 * store.go
 * $Id$
 */

package agent

import (
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
)

var store storage.Storage

var mem = new(memory.MemStorage)

func init() {
	mem.Metrics = make(map[string]*string)
}

func Storage() storage.Storage {
	if store == nil {
		store = mem
	}
	return store
}
