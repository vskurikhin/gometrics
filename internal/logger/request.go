/*
 * This file was last modified at 2024-02-24 17:37 by Victor N. Skurikhin.
 * request.go
 * $Id$
 */

package logger

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		Log.Debug("got incoming HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.String("duration", fmt.Sprintf("%v", duration)),
		)

		Log.Debug("say outgoing HTTP response",
			zap.String("status", r.Method),
			zap.String("size", fmt.Sprintf("%d", responseData.size)),
		)
	}
	return http.HandlerFunc(logFn)
}
