/*
 * This file was last modified at 2024-02-29 23:37 by Victor N. Skurikhin.
 * logging.go
 * $Id$
 */

package logger

import (
	"github.com/vskurikhin/gometrics/internal/chimiddleware"
	"github.com/vskurikhin/gometrics/internal/util"
	"net/http"
	"time"
)

// Logging returns a logger handler.
func Logging(next http.Handler) http.Handler {
	logFunc := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {

			iZapFields := util.MakeZapFields()
			iZapFields.Append("uri", r.RequestURI)
			iZapFields.Append("method", r.Method)
			iZapFields.Append("duration", time.Since(start))
			Log.Debug("got incoming HTTP request", iZapFields...)

			oZapFields := util.MakeZapFields()
			oZapFields.AppendInt("status", ww.Status())
			oZapFields.AppendInt("size", ww.BytesWritten())
			Log.Debug("say outgoing HTTP response", oZapFields...)
		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(logFunc)
}
