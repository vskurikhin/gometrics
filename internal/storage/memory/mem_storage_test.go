/*
 * This file was last modified at 2024-02-11 00:09 by Victor N. Skurikhin.
 * mem_storage_test.go
 * $Id$
 */

package memory

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemStorage(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{"Test String() method of Alloc with type Number", "Alloc", "Alloc"},
		{"Test String() method of RandomValue with type Number", "RandomValue", "RandomValue"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := Instance()
			storage.Put(test.name, &test.input)
			got := storage.Get(test.name)
			assert.Equal(t, test.want, *got)
		})
	}
}
