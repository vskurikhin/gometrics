/*
 * This file was last modified at 2024-02-11 00:39 by Victor N. Skurikhin.
 * metric_name.go
 * $Id$
 */

package types

import (
	"slices"
	"strings"
)

type Name int

const (
	_none Name = iota
	Alloc
	BuckHashSys
	Frees
	GCCPUFraction
	GCSys
	HeapAlloc
	HeapIdle
	HeapInuse
	HeapObjects
	HeapReleased
	HeapSys
	LastGC
	Lookups
	MCacheInuse
	MCacheSys
	MSpanInuse
	MSpanSys
	Mallocs
	NextGC
	NumForcedGC
	NumGC
	OtherSys
	PauseTotalNs
	StackInuse
	StackSys
	Sys
	TotalAlloc
	PollCount
	RandomValue
)

var lowerCase []*string

func (n Name) String() string {
	return Metrics[n].name
}

func (n Name) URLPath() string {
	return Metrics[n].path
}

func Lookup(s string) Name {
	i := lookup(s)
	if i < 0 {
		return 0
	}
	return Name(i)
}

func lookup(s string) int {
	return slices.IndexFunc(lowerCase, func(e *string) bool {
		return *e == strings.ToLower(s)
	})
}
