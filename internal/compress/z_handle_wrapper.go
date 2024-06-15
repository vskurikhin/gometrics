/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * z_handle_wrapper.go
 * $Id$
 */

// Package compress сжатие
package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/vskurikhin/gometrics/internal/crypto"
	"github.com/vskurikhin/gometrics/internal/env"
	"io"
	"net/http"
	"strings"

	"github.com/vskurikhin/gometrics/internal/logger"
)

func ZHandleWrapper(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request) int) int {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {

		logger.Log.Debug("got incoming HTTP request with Content-Encoding gzip in ZHandleWrapper")
		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return http.StatusInternalServerError
		}
		cfg := env.GetServerConfig()
		crypt := crypto.GetServerCrypto(cfg)
		reader := tryDecryptRSA(crypt, body)

		// создаём *gzip.Reader, который будет читать тело запроса
		// и распаковывать его
		gz, err := gzip.NewReader(reader)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return http.StatusInternalServerError
		}
		// закрытие gzip-читателя опционально, т.к. все данные уже прочитаны и
		// текущая реализация не требует закрытия, тем не менее лучше это делать -
		// некоторые реализации могут рассчитывать на закрытие читателя
		// gz.Close() не вызывает закрытия r.Body - это будет сделано позже, http-сервером
		//nolint:multichecker,errcheck
		defer func() { _ = gz.Close() }()
		r.Body = gz
	}
	return handler(w, r)
}

func tryDecryptRSA(crypt crypto.Crypto, b []byte) io.Reader {

	if buf, ok := crypt.TryDecryptRSA(b); ok {
		return bytes.NewBuffer(buf)
	} else {
		return bytes.NewBuffer(b)
	}
}
