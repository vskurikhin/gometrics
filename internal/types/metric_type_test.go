/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * metric_type_test.go
 * $Id$
 */

package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypesString(t *testing.T) {
	var tests = []struct {
		name  string
		input Types
		want  string
	}{
		{"Test String() method of Alloc with type Types", GAUGE, "Gauge"},
		{"Test String() method of RandomValue with type Types", COUNTER, "Counter"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestTypesURLPath(t *testing.T) {
	var tests = []struct {
		name  string
		input Types
		want  string
	}{
		{"Test URLPath() method of GAUGE with type Types", GAUGE, "gauge"},
		{"Test URLPath() method of COUNTER with type Types", COUNTER, "counter"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.URLPath()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestTypesEq(t *testing.T) {
	var tests = []struct {
		name  string
		input Types
		want  string
	}{
		{"Test URLPath() method of GAUGE with type Types", GAUGE, "GAUGE"},
		{"Test URLPath() method of COUNTER with type Types", COUNTER, "Counter"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.True(t, test.input.Eq(test.want))
		})
	}
}

func TestTypesParseValue(t *testing.T) {
	var tests = []struct {
		name  string
		input Types
		want  float64
	}{
		{"Test URLPath() method of GAUGE with type Types", GAUGE, 13},
		{"Test URLPath() method of COUNTER with type Types", COUNTER, 13},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.input.ParseValue("13")
			switch f := got.(type) {
			case float64:
				assert.False(t, math.Abs(test.want-f) > math.SmallestNonzeroFloat64)
			case int:
				assert.False(t, math.Abs(test.want-float64(f)) > math.SmallestNonzeroFloat64)
			default:
				assert.Fail(t, "unknown type")
			}
			assert.Nil(t, err)
		})
	}
	t.Run("negative case", func(t *testing.T) {
		got, err := Types(math.MaxInt32).ParseValue("a")
		switch f := got.(type) {
		case int:
			assert.Fail(t, "this is not int type")
		case float64:
			assert.True(t, math.Abs(f-0.0) < math.SmallestNonzeroFloat64)
		}
		assert.NotNil(t, err)
	})
}
