/*
 * This file was last modified at 2024-04-03 07:24 by Victor N. Skurikhin.
 * main_test.go.go
 * $Id$
 */

package main

import (
	"github.com/vskurikhin/gometrics/internal/agent"
	"github.com/vskurikhin/gometrics/internal/env"
	t "github.com/vskurikhin/gometrics/internal/types"
)

var enabled = []t.Name{t.Alloc, t.BuckHashSys, t.Frees, t.GCCPUFraction, t.GCSys,
	t.HeapAlloc, t.HeapIdle, t.HeapInuse, t.HeapObjects, t.HeapReleased, t.HeapSys,
	t.LastGC, t.Lookups, t.MCacheInuse, t.MCacheSys, t.MSpanInuse, t.MSpanSys, t.Mallocs,
	t.NextGC, t.NumForcedGC, t.NumGC, t.OtherSys, t.PauseTotalNs, t.StackInuse, t.StackSys, t.Sys,
	t.TotalAlloc, t.PollCount, t.RandomValue}

func main() {

	env.InitAgent()
	agent.Storage()

	go agent.Poll(enabled)
	agent.Reports(enabled)
}
