/*
 * This file was last modified at 2024-02-10 12:28 by Victor N. Skurikhin.
 * metrics_test.go
 * $Id$
 */

package types

import (
	"github.com/stretchr/testify/assert"
	"math"
	"runtime"
	"testing"
)

func TestMetricGetMetric(t *testing.T) {
	var tests = []struct {
		name  string
		input Name
		want  *metric
	}{
		{"Test GetMetric() method of Alloc with type Number", Alloc, &Metrics[Alloc]},
		{"Test GetMetric() method of RandomValue with type Number", RandomValue, &Metrics[RandomValue]},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.GetMetric()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetricString(t *testing.T) {
	var tests = []struct {
		name  string
		input metric
		want  string
	}{
		{"Test String() method of metric[Alloc] with type metric", Metrics[Alloc], "Alloc"},
		{"Test String() method of metric[RandomValue] with type metric", Metrics[RandomValue], "RandomValue"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.String()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetricURLPath(t *testing.T) {
	var tests = []struct {
		name  string
		input metric
		want  string
	}{
		{"Test URLPath() method of metric[Alloc] with type metric", Metrics[Alloc], "alloc"},
		{"Test URLPath() method of metric[RandomValue] with type metric", Metrics[RandomValue], "randomvalue"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.URLPath()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetricTypeString(t *testing.T) {
	var tests = []struct {
		name  string
		input metric
		want  Types
	}{
		{"Test String() method of metric[Alloc] with type metric", Metrics[Alloc], GAUGE},
		{"Test String() method of metric[PollCount] with type metric", Metrics[PollCount], COUNTER},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.MetricType()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestMetrics(t *testing.T) {
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	for _, test := range Metrics {
		t.Run(test.name, func(t *testing.T) {
			switch test.Type().(type) {
			case uint64:
				value := test.FuncUint64()(memStats)
				assert.True(t, value < math.MaxUint64)
			case uint32:
				value := test.FuncUint32()(memStats)
				assert.True(t, value < math.MaxUint32)
			case float64:
				value := test.FuncFloat64()(memStats)
				assert.True(t, value > -1.0)
			}
		})
	}
}
