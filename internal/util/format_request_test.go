/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * format_request_test.go
 * $Id$
 */

package util

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/vskurikhin/gotool"
)

func TestFormatRequest(t *testing.T) {
	headers := map[string][]string{"Content-Type": {"text/plain; charset=utf-8"}}
	var tests = []struct {
		input  string
		method string
		header map[string][]string
		want   string
	}{
		{
			input:  "",
			method: http.MethodGet,
			header: headers,
			want:   "GET  \nHost: \ncontent-type: text/plain; charset=utf-8",
		},
		{
			input:  "a",
			method: http.MethodPost,
			header: headers,
			want:   "POST a \nHost: \ncontent-type: text/plain; charset=utf-8\n\n\n",
		},
	}
	for _, test := range tests {
		u := url.URL{Path: test.input}
		request := http.Request{Method: test.method, URL: &u, Header: test.header}
		got := FormatRequest(&request)
		gotool.AssertEqual(t, got, test.want)
	}
}
