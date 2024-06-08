/*
 * This file was last modified at 2024-05-28 21:57 by Victor N. Skurikhin.
 * handlers_test.go
 * $Id$
 */

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"net/http"
	"net/http/httptest"
)

func ExampleUpdateHandler() {
	w := httptest.NewRecorder()
	target := fmt.Sprintf("%s%s/%s/%s", env.UpdateURL, "counter", "Alloc", "1")
	r := httptest.NewRequest(http.MethodPost, target, nil)

	ctx := chi.NewRouteContext()

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	handlers.UpdateHandler(w, r)
}

func ExampleUpdateJSONHandler() {

	input := map[string]interface{}{
		"id":    "Alloc",
		"type":  "gauge",
		"value": 1.1,
	}

	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, env.UpdateURL, bytes.NewReader(body))

	ctx := chi.NewRouteContext()

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	handlers.UpdateJSONHandler(w, r)
}

func ExampleValueJSONHandler() {
	input := map[string]interface{}{
		"id":    "Alloc",
		"type":  "gauge",
		"value": 1.1,
	}
	body, _ := json.Marshal(input)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, env.UpdateURL, bytes.NewReader(body))

	ctx := chi.NewRouteContext()

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	handlers.ValueJSONHandler(w, r)
}
