/*
 * This file was last modified at 2024-02-11 00:39 by Victor N. Skurikhin.
 * poll.go
 * $Id$
 */

package agent

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	types2 "github.com/vskurikhin/gometrics/internal/types"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

var cnt = atomic.Uint64{}

func Poll(enabled []types2.Name) {

	memStats := new(runtime.MemStats)
	storage := memory.Instance()
	for {
		runtime.ReadMemStats(memStats)
		for _, i := range enabled {
			putSample(storage, memStats, i)
			putCustom(storage, i)
		}
		time.Sleep(env.Agent.PollInterval() * time.Second)
	}
}

func putSample(storage *memory.MemStorage, memStats *runtime.MemStats, n types2.Name) {

	metric := n.GetMetric()
	name := metric.String()
	switch metric.Type().(type) {
	case uint64:
		value := fmt.Sprintf("%d", types2.Metrics[n].FuncUint64()(memStats))
		storage.Put(name, &value)
	case uint32:
		value := fmt.Sprintf("%d", types2.Metrics[n].FuncUint32()(memStats))
		storage.Put(name, &value)
	case float64:
		value := fmt.Sprintf("%f", types2.Metrics[n].FuncFloat64()(memStats))
		storage.Put(name, &value)
	}
}

func putCustom(storage *memory.MemStorage, n types2.Name) {

	metric := n.GetMetric()
	name := metric.String()
	switch n {
	case types2.PollCount:
		value := fmt.Sprintf("%d", cnt.Add(1))
		storage.Put(name, &value)
	case types2.RandomValue:
		value := fmt.Sprintf("%d", rand.Int())
		storage.Put(name, &value)
	}
}
