/*
 * This file was last modified at 2024-06-10 18:56 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/parser"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/types"
	"net/http"
)

// UpdateHandler обработчик сбора метрик и алертинга.
// данные в формате:
//
//	POST http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
//
// :
//
//   - принимает метрики по протоколу HTTP методом POST;
//
//   - при успешном приёме возвращать http.StatusOK;
//
//   - при попытке передать запрос без имени метрики возвращать http.StatusNotFound;
//
//   - при попытке передать запрос с некорректным типом метрики или значением возвращать http.StatusBadRequest.
//
// Принимает и хранит произвольные метрики двух типов:
//
// • gauge, float64 — новое значение должно замещать предыдущее;
//
// • counter, int64 — новое значение должно добавляться к предыдущему, если какое-то значение уже было известно серверу.
func UpdateHandler(response http.ResponseWriter, request *http.Request) {

	store = server.Storage(env.GetServerConfig())

	defer func() {
		if p := recover(); p != nil {
			response.WriteHeader(http.StatusNotFound)
		}
	}()

	parsed, err := parser.Parse(request)
	if err != nil || parsed.Value() == nil {
		response.WriteHeader(parsed.Status())
		return
	}

	name := parsed.String()
	switch parsed.Type() {
	case types.COUNTER:
		store.PutCounter(name, parsed.CalcValue(store.GetCounter(name)))
	case types.GAUGE:
		store.PutGauge(name, parsed.CalcValue(store.GetGauge(name)))
	}
	response.WriteHeader(http.StatusOK)
}
