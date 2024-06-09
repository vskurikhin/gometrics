/*
 * This file was last modified at 2024-06-10 22:30 by Victor N. Skurikhin.
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
	//nolint:multichecker,errcheck
	defer func() { _ = Log.Sync() }()
}
