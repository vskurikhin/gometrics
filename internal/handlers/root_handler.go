/*
 * This file was last modified at 2024-03-01 21:40 by Victor N. Skurikhin.
 * root_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func RootHandler(response http.ResponseWriter, request *http.Request) {
	compress.ZHandleWrapper(response, request, plainRootHandler)
}

func plainRootHandler(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(root(response))
}

func root(response http.ResponseWriter) (status int) {

	response.Header().Set("Content-Type", "text/html")

	defer func() {
		if p := recover(); p != nil {

			logger.Log.Debug("func UpdateJSONHandler", zap.String("error", fmt.Sprintf("%v", p)))
			status = http.StatusNotFound
		}
	}()
	_, _ = response.Write([]byte("<html></html>"))

	return http.StatusOK
}
