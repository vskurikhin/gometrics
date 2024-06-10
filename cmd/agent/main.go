/*
 * This file was last modified at 2024-06-10 18:44 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/agent"
	"github.com/vskurikhin/gometrics/internal/env"
	t "github.com/vskurikhin/gometrics/internal/types"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
	enabled      = []t.Name{t.Alloc, t.BuckHashSys, t.Frees, t.GCCPUFraction, t.GCSys,
		t.HeapAlloc, t.HeapIdle, t.HeapInuse, t.HeapObjects, t.HeapReleased, t.HeapSys,
		t.LastGC, t.Lookups, t.MCacheInuse, t.MCacheSys, t.MSpanInuse, t.MSpanSys, t.Mallocs,
		t.NextGC, t.NumForcedGC, t.NumGC, t.OtherSys, t.PauseTotalNs, t.StackInuse, t.StackSys, t.Sys,
		t.TotalAlloc, t.PollCount, t.RandomValue}
)

func main() {

	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	cfg := env.GetAgentConfig()
	fmt.Print(cfg)
	agent.Storage()

	go agent.Poll(cfg, enabled)
	agent.Reports(cfg, enabled)
}
