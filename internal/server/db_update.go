/*
 * This file was last modified at 2024-03-18 16:24 by Victor N. Skurikhin.
 * db_update.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
	"time"
)

func DBUpdate() {
	for {
		dbUpdate()
	}
}

func dbUpdate() {
	time.Sleep(env.Server.StoreInterval() * time.Second)
	store.SaveToFile(env.Server.FileStoragePath())
}
