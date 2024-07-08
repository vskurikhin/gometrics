/*
 * This file was last modified at 2024-07-07 11:44 by Victor N. Skurikhin.
 * read.go
 * $Id$
 */

package server

import (
	"github.com/vskurikhin/gometrics/internal/env"
)

func Read(cfg env.Config) {
	if cfg.Restore() && cfg.FileStoragePath() != "" {
		cfg.
			Property().
			Storage().
			ReadFromFile(cfg.FileStoragePath())
	}
}
