/*
 * This file was last modified at 2024-06-16 13:33 by Victor N. Skurikhin.
 * read.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
)

func Read(cfg env.Config) {
	if cfg.Restore() && cfg.FileStoragePath() != "" {
		store.ReadFromFile(cfg.FileStoragePath())
	}
}
