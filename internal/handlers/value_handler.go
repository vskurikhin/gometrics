/*
 * This file was last modified at 2024-03-18 22:49 by Victor N. Skurikhin.
 * value_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/types"
	"io"
	"net/http"
	"strings"
)

func ValueHandler(response http.ResponseWriter, request *http.Request) {

	store = server.Storage()
	typ := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")
	num := types.Lookup(name)

	if num < 1 || num.GetMetric().MetricType().Eq(typ) {
		var value *string
		if num > 0 {
			name = num.String()
		}
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
