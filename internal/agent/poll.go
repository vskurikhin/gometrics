/*
 * This file was last modified at 2024-06-25 00:17 by Victor N. Skurikhin.
 * poll.go
 * $Id$
 */

// Package agent реализация агента
package agent

import (
	"context"
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/types"
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"
)

var count = atomic.Uint64{}

func Poll(ctx context.Context, cfg env.Config, enabled []types.Name) {

	memStats := new(runtime.MemStats)
	for {
		select {
		case <-ctx.Done():
			poll(enabled, memStats)
			return
		default:
			go poll(enabled, memStats)
			time.Sleep(cfg.PollInterval())
		}
	}
}

func poll(enabled []types.Name, memStats *runtime.MemStats) {

	runtime.ReadMemStats(memStats)

	for _, i := range enabled {
		putSample(i, memStats)
		putCustom(i)
	}
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

	switch n {
	case types.PollCount:
		value := fmt.Sprintf("%d", count.Add(1))
		store.Put(name, &value)
	case types.RandomValue:
		value := fmt.Sprintf("%d", rand.Int())
		store.Put(name, &value)
	}
}
