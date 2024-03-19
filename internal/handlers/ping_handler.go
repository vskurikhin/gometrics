/*
 * This file was last modified at 2024-03-18 15:50 by Victor N. Skurikhin.
 * ping_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"go.uber.org/zap"
	"net/http"
)

var dbHealthInstance = server.DBHealthInstance()

func PingHandler(response http.ResponseWriter, request *http.Request) {
	compress.ZHandleWrapper(response, request, plainPingHandler)
}

func plainPingHandler(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(ping(response))
}

func ping(response http.ResponseWriter) (status int) {

	response.Header().Set("Content-Type", "text/plain")

	defer func() {
		if p := recover(); p != nil {

			logger.Log.Debug("func PingHandler", zap.String("error", fmt.Sprintf("%v", p)))
			status = http.StatusNotFound
		}
	}()

	if dbHealthInstance.GetStatus() {
		_, _ = response.Write([]byte("Ok"))
		return http.StatusOK
	}

	response.WriteHeader(http.StatusInternalServerError)
	_, _ = response.Write([]byte("DataBase health NOT OK!"))

	return http.StatusInternalServerError
}
