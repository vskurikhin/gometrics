/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * z_handle_wrapper_test.go
 * $Id$
 */

package compress

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestZHandleWrapper(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/type/PollCount", nil)
	r.Header.Add("Content-Encoding", "gzip")
	ctx := chi.NewRouteContext()
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	ZHandleWrapper(w, r, func(writer http.ResponseWriter, request *http.Request) int {
		return http.StatusOK
	})
}
