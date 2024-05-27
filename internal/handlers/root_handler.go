/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * root_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/logger"
)

func RootHandler(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(root(response))
}

func root(response http.ResponseWriter) (status int) {

	response.Header().Set("Content-Type", "text/html")

	defer func() {
		if p := recover(); p != nil {

			logger.Log.Debug("func RootHandler", zap.String("error", fmt.Sprintf("%v", p)))
			status = http.StatusNotFound
		}
	}()
	_, _ = response.Write([]byte("<html></html>"))

	return http.StatusOK
}
