/*
 * This file was last modified at 2024-02-10 23:44 by Victor N. Skurikhin.
 * poll.go
 * $Id$
 */

package agent

import (
	"fmt"
	"github.com/vskurikhin/gometrics/api/types"
	"github.com/vskurikhin/gometrics/cmd/env"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

var cnt = atomic.Uint64{}

func Poll(enabled []types.Name) {

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

func putSample(storage *memory.MemStorage, memStats *runtime.MemStats, n types.Name) {

	metric := n.GetMetric()
	name := metric.String()
	switch metric.Type().(type) {
	case uint64:
		value := fmt.Sprintf("%d", types.Metrics[n].FuncUint64()(memStats))
		storage.Put(name, &value)
	case uint32:
		value := fmt.Sprintf("%d", types.Metrics[n].FuncUint32()(memStats))
		storage.Put(name, &value)
	case float64:
		value := fmt.Sprintf("%f", types.Metrics[n].FuncFloat64()(memStats))
		storage.Put(name, &value)
	}
}

func putCustom(storage *memory.MemStorage, n types.Name) {

	metric := n.GetMetric()
	name := metric.String()
	switch n {
	case types.PollCount:
		value := fmt.Sprintf("%d", cnt.Add(1))
		storage.Put(name, &value)
	case types.RandomValue:
		value := fmt.Sprintf("%d", rand.Int())
		storage.Put(name, &value)
	}
}
