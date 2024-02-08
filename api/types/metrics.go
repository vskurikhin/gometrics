/*
 * This file was last modified at 2024-02-08 08:43 by Victor N. Skurikhin.
 * metric_name.go
 * $Id$
 */

package types

import (
	"fmt"
	"runtime"
	"strings"
)

type metric struct {
	name        string
	path        string
	type_       interface{}
	funcString  func(stringer fmt.Stringer) string
	funcUint32  func(*runtime.MemStats) uint32
	funcUint64  func(*runtime.MemStats) uint64
	funcFloat64 func(*runtime.MemStats) float64
}

var Metrics = [...]metric{
	{name: ""},
	{name: "Alloc", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.Alloc
	}},
	{name: "BuckHashSys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.BuckHashSys
	}},
	{name: "Frees", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.Frees
	}},
	{name: "GCCPUFraction", type_: float64(0), funcFloat64: func(ms *runtime.MemStats) float64 {
		return ms.GCCPUFraction
	}},
	{name: "GCSys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.GCSys
	}},
	{name: "HeapAlloc", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.HeapAlloc
	}},
	{name: "HeapIdle", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.HeapIdle
	}},
	{name: "HeapInuse", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.HeapIdle
	}},
	{name: "HeapObjects", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.HeapObjects
	}},
	{name: "HeapReleased", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.HeapReleased
	}},
	{name: "HeapSys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.HeapSys
	}},
	{name: "LastGC", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.LastGC
	}},
	{name: "Lookups", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.Lookups
	}},
	{name: "MCacheInuse", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.MCacheInuse
	}},
	{name: "MCacheSys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.MCacheSys
	}},
	{name: "MSpanInuse", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.MSpanInuse
	}},
	{name: "MSpanSys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.MSpanSys
	}},
	{name: "Mallocs", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.Mallocs
	}},
	{name: "NextGC", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.NextGC
	}},
	{name: "NumForcedGC", type_: uint32(0), funcUint32: func(ms *runtime.MemStats) uint32 {
		return ms.NumForcedGC
	}},
	{name: "NumGC", type_: uint32(0), funcUint32: func(ms *runtime.MemStats) uint32 {
		return ms.NumGC
	}},
	{name: "OtherSys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.OtherSys
	}},
	{name: "PauseTotalNs", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.PauseTotalNs
	}},
	{name: "StackInuse", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.StackInuse
	}},
	{name: "StackSys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.StackSys
	}},
	{name: "Sys", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.Sys
	}},
	{name: "TotalAlloc", type_: uint64(0), funcUint64: func(ms *runtime.MemStats) uint64 {
		return ms.TotalAlloc
	}},
	{name: "PollCount"},
	{name: "RandomValue"},
}

func init() {

	for i := range Metrics {
		Metrics[i].path = strings.ToLower(Metrics[i].name)
		lowerCase = append(lowerCase, &Metrics[i].path)
	}
}

func (n Name) GetMetric() *metric {
	result := Metrics[n]
	return &result
}

func (m *metric) String() string {
	return m.name
}

func (m *metric) URLPath() string {
	return m.path
}

func (m *metric) Type() interface{} {
	return m.type_
}

func (m *metric) FuncString() func(stringer fmt.Stringer) string {
	return m.funcString
}

func (m *metric) FuncUint64() func(*runtime.MemStats) uint64 {
	return m.funcUint64
}

func (m *metric) FuncUint32() func(*runtime.MemStats) uint32 {
	return m.funcUint32
}

func (m *metric) FuncFloat64() func(*runtime.MemStats) float64 {
	return m.funcFloat64
}
