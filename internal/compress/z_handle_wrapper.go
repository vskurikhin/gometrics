/*
 * This file was last modified at 2024-07-05 15:51 by Victor N. Skurikhin.
 * z_handle_wrapper.go
 * $Id$
 */

// Package compress сжатие
package compress

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"github.com/vskurikhin/gometrics/internal/crypto"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"github.com/vskurikhin/gometrics/internal/logger"
)

func ZHandleWrapper(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request) int) int {

	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {

		logger.Log.Debug("got incoming HTTP request with Content-Encoding gzip in ZHandleWrapper")

		// создаём *gzip.Reader, который будет читать тело запроса
		// и распаковывать его
		gz, err := gzip.NewReader(r.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return http.StatusInternalServerError
		}
		body, err := io.ReadAll(gz)

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

		crypt := crypto.GetServerCrypto()
		var reader io.ReadCloser
		key := r.Header.Get("X-Content-Encrypting")

		if key != "" {
			reader = tryDecryptAES(crypt, key, body)
		} else {
			reader = io.NopCloser(bytes.NewBuffer(body))
		}
		r.Body = reader
	}
	return handler(w, r)
}

func tryDecryptAES(crypt crypto.Crypto, key string, body []byte) io.ReadCloser {

	logger.Log.Debug("post DecryptRSA", zap.String("key", key))
	secretKey, err := decodeBase64AndDecryptRSA(crypt, key)

	if err == nil && len(body) > 2 {
		if buf, e := crypt.DecryptAES(secretKey, body); e != nil {
			logger.Log.Debug("func ZHandleWrapper in DecryptAES", zap.String("error", fmt.Sprintf("%v", e)))
			return io.NopCloser(bytes.NewBuffer(body))
		} else {
			logger.Log.Debug("func ZHandleWrapper in DecryptAES Ok!")
			return io.NopCloser(bytes.NewBuffer(buf))
		}
	}
	logger.Log.Debug("func ZHandleWrapper", zap.String("error", fmt.Sprintf("%v", err)))
	return io.NopCloser(bytes.NewBuffer(body))
}

func decodeBase64AndDecryptRSA(crypt crypto.Crypto, b64 string) ([]byte, error) {

	decoded, err := base64.StdEncoding.DecodeString(b64)

	if err != nil {
		return nil, err
	}
	if plain, err := crypt.DecryptRSA(decoded); err != nil {
		logger.Log.Debug("post DecryptRSA",
			zap.String("error", fmt.Sprintf("%v", err)),
		)
		return nil, err
	} else {
		return plain, nil
	}
}
