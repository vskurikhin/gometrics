/*
 * This file was last modified at 2024-03-18 19:59 by Victor N. Skurikhin.
 * value_handler_test.go
 * $Id$
 */

package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vskurikhin/gometrics/api/names"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValueHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name     string
		input    string
		type_    string
		variable string
		want     want
	}{
		{
			name:     "positive test #0",
			input:    "0.1",
			type_:    "gauge",
			variable: "_none",
			want: want{
				code:        200,
				response:    "0.1\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:     "positive test #1",
			input:    "1.1",
			type_:    "gauge",
			variable: "Alloc",
			want: want{
				code:        200,
				response:    "1.1\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name:     "negative test #0",
			input:    "1.1",
			type_:    "counter",
			variable: "Alloc",
			want: want{
				code:        404,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, names.ValueURL+"{type}/{name}", nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("type", test.type_)
			rctx.URLParams.Add("name", test.variable)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			if test.input != "" {
				store.Put(test.variable, &test.input)
			} else {
				store.Put(test.variable, nil)
			}

			ValueHandler(w, r)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			//goland:noinspection GoUnhandledErrorResult
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
