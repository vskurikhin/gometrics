/*
 * This file was last modified at 2024-02-11 00:38 by Victor N. Skurikhin.
 * parser_test.go
 * $Id$
 */

package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gometrics/internal/types"
	"net/http"
	"net/url"
	"testing"
)

func TestParser(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  parser
	}{
		{
			"Test Parse() function for PollCount with type: parser",
			"/update/counter/pollcount/31",
			parser{
				type_:         types.COUNTER,
				number:        types.PollCount,
				name:          "PollCount",
				originalValue: "31",
				parsedValue:   int(31),
				httpStatus:    http.StatusOK,
			},
		},
		{
			"Test String() method for GCCPUFraction with type parser",
			"/update/gauge/gccpufraction/676524.874",
			parser{
				type_:         types.GAUGE,
				number:        types.GCCPUFraction,
				name:          "GCCPUFraction",
				originalValue: "676524.874",
				parsedValue:   float64(676524.874),
				httpStatus:    http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := &http.Request{
				Method: http.MethodPost,
				URL: &url.URL{
					Path: test.input,
				},
			}
			test.want.request = r
			got, err := Parse(r)
			assert.Nil(t, err)
			assert.Equal(t, test.want, *got)
		})
	}
}
