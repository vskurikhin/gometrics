/*
 * This file was last modified at 2024-05-28 21:57 by Victor N. Skurikhin.
 * updates_json_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"net/http"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/util"
)

func UpdatesJSONHandler(response http.ResponseWriter, request *http.Request) {
	if store == nil {
		store = server.Storage()
	}
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

	status, err := updatesJSON(response, request)
	if err != nil {
		return status
	}

	return http.StatusOK
}

func updatesJSON(response http.ResponseWriter, request *http.Request) (int, error) {

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
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
