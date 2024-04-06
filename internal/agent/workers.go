/*
 * This file was last modified at 2024-04-06 18:23 by Victor N. Skurikhin.
 * workers.go
 * $Id$
 */

package agent

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
	"go.uber.org/zap"
	"net/http"
)

func Workers(enabled []types.Name) chan<- int {

	jobs := make(chan int, env.Agent.RateLimit())

	for w := 1; w <= env.Agent.RateLimit(); w++ {
		go worker(enabled, w, jobs)
	}
	return jobs
}

func worker(enabled []types.Name, w int, jobs <-chan int) {

	client := http.Client{}

	for i := range jobs {
		logger.Log.Debug("jobs",
			zap.String("number", fmt.Sprintf("%d", i)),
			zap.String("worker", fmt.Sprintf("%d", w)),
		)
		reports(enabled, &client)
	}
}
