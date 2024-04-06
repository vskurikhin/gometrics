/*
 * This file was last modified at 2024-04-06 18:38 by Victor N. Skurikhin.
 * poll.go
 * $Id$
 */

package agent

import (
	"fmt"
	psmem "github.com/shirou/gopsutil/v3/mem"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/types"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

var count = atomic.Uint64{}

func Poll(enabled []types.Name, jobs chan<- int) {

	var number int
	memStats := new(runtime.MemStats)
	for {
		poll(enabled, memStats)
		jobs <- inc(&number)
	}
}

func poll(enabled []types.Name, memStats *runtime.MemStats) {

	runtime.ReadMemStats(memStats)

	for _, i := range enabled {
		putSample(i, memStats)
		putCustom(i)
	}
	time.Sleep(env.Agent.PollInterval() * time.Second)
}

func putSample(n types.Name, memStats *runtime.MemStats) {

	metric := n.GetMetric()
	name := metric.String()

	switch metric.Type().(type) {
	case uint64:
		value := fmt.Sprintf("%d", types.Metrics[n].FuncUint64()(memStats))
		store.Put(name, &value)
	case uint32:
		value := fmt.Sprintf("%d", types.Metrics[n].FuncUint32()(memStats))
		store.Put(name, &value)
	case float64:
		value := fmt.Sprintf("%f", types.Metrics[n].FuncFloat64()(memStats))
		store.Put(name, &value)
	}
}

func putCustom(n types.Name) {

	metric := n.GetMetric()
	name := metric.String()
	v, _ := psmem.VirtualMemory()

	switch n {
	case types.PollCount:
		value := fmt.Sprintf("%d", count.Add(1))
		store.Put(name, &value)
	case types.RandomValue:
		value := fmt.Sprintf("%d", rand.Int())
		store.Put(name, &value)
	case types.TotalMemory:
		value := fmt.Sprintf("%d", v.Total)
		store.Put(name, &value)
	case types.FreeMemory:
		value := fmt.Sprintf("%d", v.Free)
		store.Put(name, &value)
	case types.CPUutilization1:
		value := fmt.Sprintf("%f", v.UsedPercent)
		store.Put(name, &value)
	}
}

func inc(i *int) int {
	*i++
	return *i
}
