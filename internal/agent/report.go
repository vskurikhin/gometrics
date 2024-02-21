/*
 * This file was last modified at 2024-02-12 21:13 by Victor N. Skurikhin.
 * report.go
 * $Id$
 */

package agent

import (
	"fmt"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"github.com/vskurikhin/gometrics/internal/types"
	"net/http"
	"net/url"
	"os"
	"time"
)

func Report(enabled []types.Name) {

	client := http.Client{}
	memStorage := memory.Instance()
	for {
		report(enabled, memStorage, &client)
	}
}

func report(enabled []types.Name, storage storage.Storage, client *http.Client) {

	time.Sleep(env.Agent.ReportInterval() * time.Second)

	for _, i := range enabled {
		post(i, client, storage)
	}
}

func post(n types.Name, client *http.Client, storage storage.Storage) {

	metric := n.GetMetric()
	name := metric.String()
	value := storage.Get(name)

	if value != nil {
		path := urlPrintf(metric.MetricType(), metric, part(*value))
		request, err := http.NewRequest(http.MethodPost, path, nil)
		if err != nil {
			panic(err)
		}
		postDo(client, request)
	}
}

func postDo(client *http.Client, request *http.Request) {

	defer func() {
		if p := recover(); p != nil {
			switch p.(type) {
			case *url.Error:
			default:
				//goland:noinspection GoUnhandledErrorResult
				fmt.Fprintf(os.Stderr, "post error: %T", p)
			}
		}
	}()

	request.Header.Add("Content-Type", "text/plain")
	response, err := client.Do(request)
	defer func() {
		if response != nil {
			//goland:noinspection GoUnhandledErrorResult
			response.Body.Close()
		}
	}()

	if err != nil {
		panic(err)
	}
}

func urlPrintf(parts ...urlPart) string {

	path := *env.Agent.URLHost() + names.Update

	for _, part := range parts {
		path += "/" + part.URLPath()
	}
	return path
}
