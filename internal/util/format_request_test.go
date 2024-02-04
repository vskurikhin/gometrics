/*
 * This file was last modified at 2024-02-04 17:13 by Victor N. Skurikhin.
 * format_request_test.go
 * $Id$
 */

package util

import (
	"github.com/vskurikhin/gotool"
	"net/http"
	"net/url"
	"testing"
)

func TestFormatRequest(t *testing.T) {
	var tests = []struct {
		input string
		want  string
	}{
		{"", "  \nHost: "},
		{"a", " a \nHost: "},
	}
	for _, test := range tests {
		u := url.URL{Path: test.input}
		request := http.Request{URL: &u}
		got := FormatRequest(&request)
		gotool.AssertEqual(t, got, test.want)
	}
}
