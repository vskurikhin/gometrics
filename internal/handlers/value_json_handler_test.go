/*
 * This file was last modified at 2024-06-10 17:45 by Victor N. Skurikhin.
 * value_json_handler_test.go
 * $Id$
 */

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/server"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValueJSONHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name     string
		input    map[string]interface{}
		type_    string
		variable string
		want     want
	}{
		{
			name: "positive test #1",
			input: map[string]interface{}{
				"id":    "Alloc",
				"type":  "gauge",
				"value": 1.1,
			},
			type_:    "gauge",
			variable: "Alloc",
			want: want{
				code:        200,
				response:    "{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":1.1}",
				contentType: "application/json",
			},
		},
		{
			name: "negative test #1",
			input: map[string]interface{}{
				"id":    "Alloc",
				"type":  "gauge",
				"value": "1.1",
			},
			type_:    "",
			variable: "Alloc",
			want: want{
				code:        404,
				response:    "",
				contentType: "application/json",
			},
		},
	}
	cfg := env.GetServerConfig()
	store = server.Storage(cfg)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			body, _ := json.Marshal(test.input)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, env.UpdateURL, bytes.NewReader(body))

			ctx := chi.NewRouteContext()

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			ValueJSONHandler(w, r)

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
