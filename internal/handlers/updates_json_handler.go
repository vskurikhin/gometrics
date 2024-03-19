/*
 * This file was last modified at 2024-03-19 09:31 by Victor N. Skurikhin.
 * updates_json_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/util"
	"go.uber.org/zap"
	"net/http"
)

func UpdatesJSONHandler(response http.ResponseWriter, request *http.Request) {
	store = server.Storage()
	compress.ZHandleWrapper(response, request, plainUpdatesJSONHandler)
}

func plainUpdatesJSONHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(updatesJSONHandler(response, request))
}

func updatesJSONHandler(response http.ResponseWriter, request *http.Request) (status int) {

	response.Header().Set("Content-Type", "application/json")

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func UpdatesJSONHandler", zap.String("error", fmt.Sprintf("%v", p)))
			status = http.StatusNotFound
		}
	}()
	updatesJSON(response, request)

	return http.StatusOK
}

func updatesJSON(response http.ResponseWriter, request *http.Request) {

	metrics := make(dto.Metrics, 0)

	if err := easyjson.UnmarshalFromReader(request.Body, &metrics); err != nil {
		panic(err)
	}

	for _, metric := range metrics {
		zapFields := util.ZapFieldsMetric(&metric)
		logger.Log.Debug("got incoming HTTP request with JSON in updatesJSON", zapFields.Slice()...)
	}
	store.PutSlice(metrics)

	if _, err := easyjson.MarshalToWriter(metrics, response); err != nil {
		panic(err)
	}
}
