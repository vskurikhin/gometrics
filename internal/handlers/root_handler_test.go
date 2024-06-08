/*
 * This file was last modified at 2024-05-28 21:57 by Victor N. Skurikhin.
 * root_handler_test.go
 * $Id$
 */

package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/type/PollCount", nil)

	ctx := chi.NewRouteContext()

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	RootHandler(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
