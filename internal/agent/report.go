/*
 * This file was last modified at 2024-02-07 21:57 by Victor N. Skurikhin.
 * report.go
 * $Id$
 */

package agent

import (
	"fmt"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/api/types"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"net/http"
	"os"
	"time"
)

func Report(enabled []types.Name) {

	client := http.Client{
		Timeout: time.Second * 1, // интервал ожидания: 1 секунда
	}
	storage := memory.Instance()
	for {
		for _, i := range enabled {
			post(i, storage, &client)
		}
		time.Sleep(names.ReportInterval * time.Second)
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
		err = response.Body.Close()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "post error: %v", err)
		}
	}()
	if err != nil {
		panic(err)
	}
}

func urlPrintf(parts ...urlPart) string {

	path := names.UpdateURLClient
	for _, part := range parts {
		path += "/" + part.URLPath()
	}
	return path
}
