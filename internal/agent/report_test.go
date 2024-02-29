/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * report_test.go
 * $Id$
 */

package agent

import (
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReport(t *testing.T) {

	env.InitAgent()
	server := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}),
	)
	defer server.Close()

	client := server.Client()

	var tests = []struct {
		name   string
		input1 []types.Name
		input2 string
		want   string
	}{
		{name: "positive test #0", input1: []types.Name{types.Alloc}, input2: "0", want: "Alloc"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			report(test.input1, client)
		})
	}
}
