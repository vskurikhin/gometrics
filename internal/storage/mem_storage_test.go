/*
 * This file was last modified at 2024-05-28 21:57 by Victor N. Skurikhin.
 * mem_storage_test.go
 * $Id$
 */

package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gometrics/internal/dto"
	"math/rand"
	"testing"
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
