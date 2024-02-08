/*
 * This file was last modified at 2024-02-04 12:18 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	t "github.com/vskurikhin/gometrics/api/types"
	"github.com/vskurikhin/gometrics/internal/agent"
)

var enabled = []t.Name{t.Alloc, t.BuckHashSys, t.Frees, t.GCCPUFraction, t.GCSys,
	t.HeapAlloc, t.HeapIdle, t.HeapInuse, t.HeapObjects, t.HeapReleased, t.HeapSys,
	t.LastGC, t.Lookups, t.MCacheInuse, t.MCacheSys, t.MSpanInuse, t.MSpanSys, t.Mallocs,
	t.NextGC, t.NumForcedGC, t.NumGC, t.OtherSys, t.PauseTotalNs, t.StackInuse, t.StackSys, t.Sys,
	t.TotalAlloc, t.PollCount, t.RandomValue}

func main() {

	go agent.Poll(enabled)
	agent.Report(enabled)
}
