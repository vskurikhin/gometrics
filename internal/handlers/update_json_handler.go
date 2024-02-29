/*
 * This file was last modified at 2024-02-29 23:37 by Victor N. Skurikhin.
 * update_json_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
	"github.com/vskurikhin/gometrics/internal/util"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

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
	updateJSON(response, request)

	return http.StatusOK
}

func updateJSON(response http.ResponseWriter, request *http.Request) {

	metric := dto.Metrics{}

	if err := easyjson.UnmarshalFromReader(request.Body, &metric); err != nil {
		panic(err)
	}

	zapFields := util.MakeZapFields()
	zapFields.Append("metric", metric)

	if metric.Value != nil {
		zapFields.AppendFloat("Value", *metric.Value)
	}
	if metric.Delta != nil {
		zapFields.AppendInt("Delta", *metric.Delta)
	}
	logger.Log.Debug("got incoming HTTP request with JSON in updateJSON", zapFields...)
	updateMetric(&metric)

	if _, err := easyjson.MarshalToWriter(metric, response); err != nil {
		panic(err)
	}
}

func updateMetric(metric *dto.Metrics) {

	name := strings.ToLower(metric.ID)
	value := store.Get(strings.ToLower(metric.ID))

	switch {
	case types.GAUGE.Eq(metric.MType):
		value := fmt.Sprintf("%.12f", *metric.Value)
		store.Put(name, &value)
	case types.COUNTER.Eq(metric.MType):
		*metric.Delta = calcMetricDelta(metric, value)
		value := fmt.Sprintf("%d", *metric.Delta)
		store.Put(name, &value)
	}
}

func calcMetricDelta(metric *dto.Metrics, value *string) int64 {

	var err error
	var i64 int64

	if value != nil {
		i64, err = strconv.ParseInt(*value, 10, 64)
	}
	if err != nil {
		panic(err)
	}
	if metric.Delta != nil {
		i64 += *metric.Delta
	}
	return i64
}
