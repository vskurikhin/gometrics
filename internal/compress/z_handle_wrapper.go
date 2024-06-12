/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * z_handle_wrapper.go
 * $Id$
 */

// Package compress сжатие
package compress

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/vskurikhin/gometrics/internal/logger"
)

func ZHandleWrapper(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request)) {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		logger.Log.Debug("got incoming HTTP request with Content-Encoding gzip in ZHandleWrapper")
		// создаём *gzip.Reader, который будет читать тело запроса
		// и распаковывать его
		gz, err := gzip.NewReader(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// закрытие gzip-читателя опционально, т.к. все данные уже прочитаны и
		// текущая реализация не требует закрытия, тем не менее лучше это делать -
		// некоторые реализации могут рассчитывать на закрытие читателя
		// gz.Close() не вызывает закрытия r.Body - это будет сделано позже, http-сервером
		//nolint:multichecker,errcheck
		defer func() {
			_ = gz.Close()
		}()
		r.Body = gz
	}
	handler(w, r)
}
