/*
 * This file was last modified at 2024-02-04 13:29 by Victor N. Skurikhin.
 * update_handler.go
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

func ValueHandler(w http.ResponseWriter, r *http.Request) {

	typ := chi.URLParam(r, "type")
	name := chi.URLParam(r, "name")
	n := types.Lookup(name)

	if n < 1 || n.GetMetric().MetricType().Eq(typ) {
		var value *string
		if n > 0 {
			value = storage.Get(n.String())
		} else {
			value = storage.Get(name)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if value == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, err := io.WriteString(w, fmt.Sprintf("%s\n", *value))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
