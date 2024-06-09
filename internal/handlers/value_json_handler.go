/*
 * This file was last modified at 2024-06-11 10:35 by Victor N. Skurikhin.
 * value_json_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/server"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
	"github.com/vskurikhin/gometrics/internal/util"
)

// ValueJSONHandler обработчик сбора метрик и алертинга, получения метрик с сервера.
//
//		POST value/
//	 Content-Type: application/json
//
// Обмен с сервером организуйте с использованием следующей структуры:
//
//	type Metrics struct {
//	    ID    string   `json:"id"`              // имя метрики
//	    MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
//	    Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
//	    Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
//	}
func ValueJSONHandler(response http.ResponseWriter, request *http.Request) {
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

	metric := dto.Metric{}

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

func valueMetric(metric *dto.Metric) {

	var err error
	var name string
	num := types.Lookup(metric.ID)

	if num > 0 {
		name = num.String()
	} else {
		name = metric.ID
	}
	store = server.Storage(env.GetServerConfig())

	switch {
	case types.GAUGE.Eq(metric.MType):
		value := store.GetGauge(name)
		metric.Value = new(float64)
		if value != nil {
			*metric.Value, err = strconv.ParseFloat(*value, 64)
		}
	case types.COUNTER.Eq(metric.MType):
		value := store.GetCounter(name)
		metric.Delta = new(int64)
		if value != nil {
			*metric.Delta, err = strconv.ParseInt(*value, 10, 64)
		}
	}
	util.IfErrorThenPanic(err)
}
