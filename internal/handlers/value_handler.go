/*
 * This file was last modified at 2024-02-10 15:27 by Victor N. Skurikhin.
 * value_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/api/types"
	"io"
	"net/http"
)

func ValueHandler(response http.ResponseWriter, request *http.Request) {

	typ := chi.URLParam(request, "type")
	name := chi.URLParam(request, "name")
	num := types.Lookup(name)

	if num < 1 || num.GetMetric().MetricType().Eq(typ) {
		var value *string
		if num > 0 {
			value = storage.Get(num.String())
		} else {
			value = storage.Get(name)
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
