/*
 * This file was last modified at 2024-03-02 19:58 by Victor N. Skurikhin.
 * read.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
)

func Read() {
	if env.Server.Restore() {
		store.ReadFromFile(env.Server.FileStoragePath())
	}
}
