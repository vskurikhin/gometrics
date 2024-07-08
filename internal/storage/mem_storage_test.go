/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * mem_storage_test.go
 * $Id$
 */

package storage

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vskurikhin/gometrics/internal/dto"
)

func TestMemStorage(t *testing.T) {

	test := "test"

	var mem = new(MemStorage)
	mem.Metrics = make(map[string]*string)

	mem.Put("test", &test)
	assert.Equal(t, &test, mem.Get("test"))

	mem.PutCounter("test", &test)
	assert.Equal(t, &test, mem.GetCounter("test"))

	mem.PutGauge("test", &test)
	assert.Equal(t, &test, mem.GetGauge("test"))

	metric := new(dto.Metric)
	metrics := make(dto.Metrics, 0)
	metrics = append(metrics, *metric)

	mem.PutSlice(metrics)

	fileName := fmt.Sprintf("test%d.txt", rand.Uint32())

	mem.SaveToFile(fileName)
	mem.ReadFromFile(fileName)
}
