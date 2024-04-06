/*
 * This file was last modified at 2024-04-06 18:34 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/agent"
	"github.com/vskurikhin/gometrics/internal/env"
	t "github.com/vskurikhin/gometrics/internal/types"
	"time"
)

var enabled = []t.Name{t.Alloc, t.BuckHashSys, t.Frees, t.GCCPUFraction, t.GCSys,
	t.HeapAlloc, t.HeapIdle, t.HeapInuse, t.HeapObjects, t.HeapReleased, t.HeapSys,
	t.LastGC, t.Lookups, t.MCacheInuse, t.MCacheSys, t.MSpanInuse, t.MSpanSys, t.Mallocs,
	t.NextGC, t.NumForcedGC, t.NumGC, t.OtherSys, t.PauseTotalNs, t.StackInuse, t.StackSys, t.Sys,
	t.TotalAlloc, t.PollCount, t.RandomValue, t.TotalMemory, t.FreeMemory, t.CPUutilization1}

func worker(id int, jobs <-chan int) {
	for j := range jobs {
		// для наглядности будем выводить какой рабочий начал работу и кго задачу
		fmt.Println("рабочий", id, "запущен задача", j)
		// немного замедлим выполнение рабочего
		time.Sleep(time.Second)
		// для наглядности выводим какой рабочий завершил какую задачу
		fmt.Println("рабочий", id, "закончил задача", j)
		// отправляем результат в канал результатов
	}
}

func main() {
	env.InitAgent()
	agent.Storage()
	agent.Poll(enabled, agent.Workers(enabled))
}
