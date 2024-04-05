/*
 * This file was last modified at 2024-04-03 08:47 by Victor N. Skurikhin.
 * update_handler_test.go
 * $Id$
 */

package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vskurikhin/gometrics/internal/env"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
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
			name:     "positive test #1",
			input:    "1.1",
			type_:    "gauge",
			variable: "Alloc",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name:     "negative test #1",
			input:    "1.1",
			type_:    "",
			variable: "Alloc",
			want: want{
				code:        400,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			target := fmt.Sprintf("%s%s/%s/%s", env.UpdateURL, test.type_, test.variable, test.input)
			r := httptest.NewRequest(http.MethodPost, target, nil)

			rctx := chi.NewRouteContext()

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			UpdateHandler(w, r)

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

func TestUpdateHandlerNegative(t *testing.T) {
	oldStorage := store
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
			name:     "negative test #0",
			input:    "1.1",
			type_:    "gauge",
			variable: "Alloc",
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			target := fmt.Sprintf("%s%s/%s/%s", env.UpdateURL, test.type_, test.variable, test.input)
			r := httptest.NewRequest(http.MethodPost, target, nil)

			rctx := chi.NewRouteContext()

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			store = nil
			UpdateHandler(w, r)

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
	store = oldStorage
}
