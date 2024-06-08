/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * parser_test.go
 * $Id$
 */

package parser

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vskurikhin/gometrics/internal/types"
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
			"Test Parse() method for GCCPUFraction with type parser",
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
		{
			"positive test #0",
			"/update/gauge/nil/676524.874",
			parser{
				type_:         types.GAUGE,
				number:        0,
				name:          "nil",
				originalValue: "676524.874",
				parsedValue:   float64(676524.874),
				httpStatus:    http.StatusOK,
			},
		},
		{
			"negative test #0",
			"",
			parser{httpStatus: http.StatusNotFound},
		},
		{
			"negative test #1",
			"/update/nil/nil/0",
			parser{httpStatus: http.StatusBadRequest},
		},
		{
			"negative test #3",
			"/update/gauge/nil/a",
			parser{httpStatus: http.StatusBadRequest},
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
			if test.want.httpStatus == http.StatusOK {
				test.want.request = r
			}
			got, err := Parse(r)
			if test.want.httpStatus == http.StatusOK {
				assert.Nil(t, err)
			}
			assert.Equal(t, test.want, *got)
		})
	}
}

func TestString(t *testing.T) {
	var tests = []struct {
		name  string
		input parser
		want  string
	}{
		{
			"Test String() function for PollCount with type: parser",
			parser{name: "PollCount"},
			"PollCount",
		},
		{
			"Test String() method for GCCPUFraction with type parser",
			parser{name: "GCCPUFraction"},
			"GCCPUFraction",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestValue(t *testing.T) {
	var tests = []struct {
		name  string
		input parser
		want  interface{}
	}{
		{
			"Test Value() function for PollCount with type: parser",
			parser{parsedValue: 31},
			31,
		},
		{
			"Test Value() method for GCCPUFraction with type parser",
			parser{parsedValue: float64(676524.874)},
			float64(676524.874),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Value()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestStatus(t *testing.T) {
	var tests = []struct {
		name  string
		input parser
		want  int
	}{
		{
			"Test Status() function for PollCount with type: parser",
			parser{httpStatus: 200},
			http.StatusOK,
		},
		{
			"Test Status() method for GCCPUFraction with type parser",
			parser{httpStatus: 404},
			http.StatusNotFound,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.Status()
			assert.Equal(t, test.want, got)
		})
	}
}

//goland:noinspection GoUnhandledErrorResult
func TestCalcValue(t *testing.T) {
	var tests = []struct {
		name  string
		input parser
		old   string
		want  string
	}{
		{
			"Test Parse() function for PollCount with type: parser",
			parser{
				type_:         types.COUNTER,
				number:        types.PollCount,
				name:          "PollCount",
				originalValue: "31",
				parsedValue:   int(31),
				httpStatus:    http.StatusOK,
			},
			"11",
			"42",
		},
		{
			"Test Parse() method for GCCPUFraction with type parser",
			parser{
				type_:         types.GAUGE,
				number:        types.GCCPUFraction,
				name:          "GCCPUFraction",
				originalValue: "676524.874",
				parsedValue:   float64(676524.874),
				httpStatus:    http.StatusOK,
			},
			"676524.741",
			"676524.874",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.CalcValue(&test.old)
			assert.Equal(t, test.want, *got)
		})
	}
	t.Run("negative test #2", func(t *testing.T) {
		p := parser{originalValue: "test"}
		got := p.CalcValue(nil)
		assert.Equal(t, "test", *got)
	})
	t.Run("negative test #3", func(t *testing.T) {
		p := parser{originalValue: "test"}
		x := "x"
		got := p.CalcValue(&x)
		fmt.Fprintf(os.Stderr, "got: %v\n", got)
		assert.Nil(t, got)
	})
}
