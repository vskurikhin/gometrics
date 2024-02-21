/*
 * This file was last modified at 2024-02-04 17:03 by Victor N. Skurikhin.
 * split_path_test.go
 * $Id$
 */

package util

import (
	"github.com/vskurikhin/gotool"
	"net/http"
	"net/url"
	"testing"
)

func TestSplitPath(t *testing.T) {
	var tests = []struct {
		input string
		want  []string
	}{
		{"", []string{""}},
		{"a", []string{"a"}},
		{"/a", []string{"a"}},
		{"/a/", []string{"a"}},
		{"/a/b", []string{"a", "b"}},
		{"/a/b/", []string{"a", "b"}},
	}
	for _, test := range tests {
		u := url.URL{Path: test.input}
		request := http.Request{URL: &u}
		got := SplitPath(&request)
		gotool.AssertEqual(t, got, test.want)
	}
}
