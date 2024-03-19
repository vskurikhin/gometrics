/*
 * This file was last modified at 2024-03-19 12:59 by Victor N. Skurikhin.
 * logging.go
 * $Id$
 */

package logger

import (
	"github.com/google/uuid"
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

			requestID := r.Header.Get("Request-id")
			if requestID == "" {
				requestID = uuid.NewString()
				r.Header.Set("Request-id", requestID)
			}

			iZapFields := util.MakeZapFields()
			iZapFields.Append("request-id", requestID)
			iZapFields.Append("uri", r.RequestURI)
			iZapFields.Append("method", r.Method)
			iZapFields.Append("duration", time.Since(start))
			Log.Debug("got incoming HTTP request", iZapFields.Slice()...)

			oZapFields := util.MakeZapFields()
			oZapFields.Append("request-id", requestID)
			oZapFields.AppendInt("status", ww.Status())
			oZapFields.AppendInt("size", ww.BytesWritten())
			Log.Debug("say outgoing HTTP response", oZapFields.Slice()...)
		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(logFunc)
}
