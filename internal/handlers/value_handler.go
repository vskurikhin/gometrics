/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * value_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/server"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/vskurikhin/gometrics/internal/types"
)

// ValueHandler обработчик сбора метрик и алертинга.
// В ответ на запрос
//
//	GET http://<АДРЕС_СЕРВЕРА>/value/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>
//
// возвращает текущее значение метрики в текстовом виде со статусом http.StatusOK.
func ValueHandler(response http.ResponseWriter, request *http.Request) {

	typ := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")
	num := types.Lookup(name)

	if num < 1 || num.GetMetric().MetricType().Eq(typ) {
		var value *string

		if num > 0 {
			name = num.String()
		}
		store = server.Storage(env.GetServerConfig())

		switch strings.ToLower(typ) {
		case "counter":
			value = store.GetCounter(name)
		case "gauge":
			value = store.GetGauge(name)
		}

		response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if value == nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := io.WriteString(response, fmt.Sprintf("%s\n", *value))
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	response.WriteHeader(http.StatusNotFound)
}
