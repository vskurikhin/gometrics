/*
 * This file was last modified at 2024-02-29 12:50 by Victor N. Skurikhin.
 * logger.go
 * $Id$
 */

package logger

import (
	"go.uber.org/zap"
)

// Log будет доступен всему коду как синглтон.
var Log *zap.Logger

//goland:noinspection GoUnhandledErrorResult
func init() {
	if false {
		logger, _ := zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any
		Log := logger.Sugar()
		_ = Log
	} else {
		Log = zap.NewExample()
		defer Log.Sync()
	}
}
