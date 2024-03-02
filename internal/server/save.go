/*
 * This file was last modified at 2024-03-02 20:04 by Victor N. Skurikhin.
 * save.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
	"time"
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
