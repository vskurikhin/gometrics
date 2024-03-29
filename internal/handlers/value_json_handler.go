/*
 * This file was last modified at 2024-03-02 19:42 by Victor N. Skurikhin.
 * value_json_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
	"github.com/vskurikhin/gometrics/internal/util"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

func ValueJSONHandler(response http.ResponseWriter, request *http.Request) {
	compress.ZHandleWrapper(response, request, plainValueJSONHandler)
}

func plainValueJSONHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(valueJSONHandler(response, request))
}

func valueJSONHandler(response http.ResponseWriter, request *http.Request) (status int) {

	response.Header().Set("Content-Type", "application/json")

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func ValueJSONHandler", zap.String("error", fmt.Sprintf("%v", p)))
			status = http.StatusNotFound
		}
	}()
	valueJSON(response, request)

	return http.StatusOK
}

func valueJSON(response http.ResponseWriter, request *http.Request) {

	metric := dto.Metrics{}

	if err := easyjson.UnmarshalFromReader(request.Body, &metric); err != nil {
		panic(err)
	}
	zapFields := util.ZapFieldsMetric(&metric)
	logger.Log.Debug("got incoming HTTP request with JSON in valueJSON", zapFields.Slice()...)
	valueMetric(&metric)

	if _, err := easyjson.MarshalToWriter(metric, response); err != nil {
		panic(err)
	}
}

func valueMetric(metric *dto.Metrics) {

	var err error
	value := store.Get(strings.ToLower(metric.ID))

	switch {
	case types.GAUGE.Eq(metric.MType):
		metric.Value = new(float64)
		if value != nil {
			*metric.Value, err = strconv.ParseFloat(*value, 64)
		}
	case types.COUNTER.Eq(metric.MType):
		metric.Delta = new(int64)
		if value != nil {
			*metric.Delta, err = strconv.ParseInt(*value, 10, 64)
		}
	}
	if err != nil {
		panic(err)
	}
}
