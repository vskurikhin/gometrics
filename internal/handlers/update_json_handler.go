/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * update_json_handler.go
 * $Id$
 */

package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/services"
	"github.com/vskurikhin/gometrics/internal/util"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	if _, err := services.GetMetricsService(env.GetServerConfig()).DTOUpdate(ctx, &metric); err != nil {
		return http.StatusNotFound, err
	}
	if _, err := easyjson.MarshalToWriter(metric, response); err != nil {
		return http.StatusNotFound, err
	}
	return http.StatusOK, nil
}
