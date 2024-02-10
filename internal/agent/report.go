/*
 * This file was last modified at 2024-02-10 23:50 by Victor N. Skurikhin.
 * report.go
 * $Id$
 */

package agent

import (
	"github.com/vskurikhin/gometrics/api/types"
	"github.com/vskurikhin/gometrics/cmd/env"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"net/http"
	"time"
)

func Report(enabled []types.Name) {

	client := http.Client{}
	storage := memory.Instance()
	for {
		time.Sleep(env.Agent.ReportInterval() * time.Second)
		for _, i := range enabled {
			post(i, storage, &client)
		}
	}
}

func post(n types.Name, storage *memory.MemStorage, client *http.Client) {

	metric := n.GetMetric()
	name := metric.String()
	value := storage.Get(name)

	if value != nil {
		path := urlPrintf(metric.MetricType(), metric, part(*value))
		request, err := http.NewRequest(http.MethodPost, path, nil)
		if err != nil {
			panic(err)
		}
		postDo(request, client)
	}
}

func postDo(request *http.Request, client *http.Client) {

	request.Header.Add("Content-Type", "text/plain")
	response, err := client.Do(request)
	defer func() {
		if response != nil {
			//goland:noinspection GoUnhandledErrorResult
			response.Body.Close()
		}
		recover()
	}()
	if err != nil {
		panic(err)
	}
}

func urlPrintf(parts ...urlPart) string {

	path := *env.Agent.URLHost()

	for _, part := range parts {
		path += "/" + part.URLPath()
	}
	return path
}
