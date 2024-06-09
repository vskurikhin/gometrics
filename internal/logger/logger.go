/*
 * This file was last modified at 2024-02-29 12:50 by Victor N. Skurikhin.
 * logger.go
 * $Id$
 */

// Package logger настройки логгирования
package logger

import (
	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
var Log *zap.Logger

//goland:noinspection GoUnhandledErrorResult
func init() {
	Log = zap.NewExample()
	defer func() {
		//nolint:multichecker,errcheck
		_ = Log.Sync()
	}()
}
