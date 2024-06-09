/*
 * This file was last modified at 2024-06-11 10:35 by Victor N. Skurikhin.
 * update_json_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/types"
	"github.com/vskurikhin/gometrics/internal/util"
	"go.uber.org/zap"
	"net/http"
)

// UpdateJSONHandler обработчик сбора метрик и алертинга, передачи метрик на сервер.
//
//		POST update/
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
func UpdateJSONHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(updateJSONHandler(response, request))
}

func updateJSONHandler(response http.ResponseWriter, request *http.Request) (status int) {

	response.Header().Set("Content-Type", "application/json")

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func UpdateJSONHandler", zap.String("error", fmt.Sprintf("%v", p)))
			status = http.StatusNotFound
		}
	}()

	status, err := updateJSON(response, request)
	if err != nil {
		return status
	}
	return http.StatusOK
}

func updateJSON(response http.ResponseWriter, request *http.Request) (int, error) {

	metric := dto.Metric{}

	if err := easyjson.UnmarshalFromReader(request.Body, &metric); err != nil {
		return http.StatusNotFound, err
	}

	zapFields := util.ZapFieldsMetric(&metric)
	logger.Log.Debug("got incoming HTTP request with JSON in updateJSON", zapFields.Slice()...)

	updateMetric(&metric)

	if _, err := easyjson.MarshalToWriter(metric, response); err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}

func updateMetric(metric *dto.Metric) {

	num := types.Lookup(metric.ID)
	var name string

	if num > 0 {
		name = num.String()
	} else {
		name = metric.ID
	}
	store = server.Storage(env.GetServerConfig())

	switch {
	case types.GAUGE.Eq(metric.MType):
		value := fmt.Sprintf("%.12f", *metric.Value)
		store.PutGauge(name, &value)
	case types.COUNTER.Eq(metric.MType):
		pv := store.GetCounter(name)
		*metric.Delta = metric.CalcDelta(pv)
		value := fmt.Sprintf("%d", *metric.Delta)
		store.PutCounter(name, &value)
	}
}
