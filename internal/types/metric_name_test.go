/*
 * This file was last modified at 2024-02-10 12:28 by Victor N. Skurikhin.
 * metric_name_test.go
 * $Id$
 */

package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNameString(t *testing.T) {
	var tests = []struct {
		name  string
		input Name
		want  string
	}{
		{"Test String() method of Alloc with type Number", Alloc, "Alloc"},
		{"Test String() method of RandomValue with type Number", RandomValue, "RandomValue"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestNameURLPath(t *testing.T) {
	var tests = []struct {
		name  string
		input Name
		want  string
	}{
		{"Test URLPath() method of Alloc with type Number", Alloc, "alloc"},
		{"Test URLPath() method of RandomValue with type Number", RandomValue, "randomvalue"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.URLPath()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestLookup(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  Name
	}{
		{"Test Lookup() function for Alloc with type Number", "alloc", Alloc},
		{"Test Lookup() function for RandomValue with type Number", "RandomValue", RandomValue},
		{"Test Lookup() function for __NONE__", "__NONE__", _none},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Lookup(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
