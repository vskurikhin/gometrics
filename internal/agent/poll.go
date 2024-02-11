/*
 * This file was last modified at 2024-02-12 20:43 by Victor N. Skurikhin.
 * poll.go
 * $Id$
 */

package agent

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"github.com/vskurikhin/gometrics/internal/types"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

var count = atomic.Uint64{}

func Poll(enabled []types.Name) {

	memStats := new(runtime.MemStats)
	memStorage := memory.Instance()
	for {
		poll(enabled, memStats, memStorage)
	}
}

func poll(enabled []types.Name, memStats *runtime.MemStats, storage storage.Storage) {

	runtime.ReadMemStats(memStats)

	for _, i := range enabled {
		putSample(i, storage, memStats)
		putCustom(i, storage)
	}
	time.Sleep(env.Agent.PollInterval() * time.Second)
}

func putSample(n types.Name, storage storage.Storage, memStats *runtime.MemStats) {

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

func putCustom(n types.Name, storage storage.Storage) {

	metric := n.GetMetric()
	name := metric.String()

	switch n {
	case types.PollCount:
		value := fmt.Sprintf("%d", count.Add(1))
		storage.Put(name, &value)
	case types.RandomValue:
		value := fmt.Sprintf("%d", rand.Int())
		storage.Put(name, &value)
	}
}
