/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * save.go
 * $Id$
 */

package server

import (
	"time"

	"github.com/vskurikhin/gometrics/internal/env"
)

func Save() {
	for {
		save()
	}
}

func save() {
	time.Sleep(env.Server.StoreInterval() * time.Second)
	store.SaveToFile(env.Server.FileStoragePath())
}
