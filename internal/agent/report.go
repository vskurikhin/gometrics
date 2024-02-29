/*
 * This file was last modified at 2024-02-29 12:50 by Victor N. Skurikhin.
 * report.go
 * $Id$
 */

package agent

import (
	"bytes"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func Report(enabled []types.Name) {

	client := http.Client{}
	for {
		report(enabled, &client)
	}
}

func report(enabled []types.Name, client *http.Client) {

	time.Sleep(env.Agent.ReportInterval() * time.Second)

	for _, i := range enabled {
		post(i, client)
	}
}

func post(n types.Name, client *http.Client) {

	name := n.GetMetric().String()
	value := store.Get(name)

	if value != nil {

		mtyp := n.
			GetMetric().
			MetricType().
			URLPath()
		path := *env.Agent.URLHost() + names.UpdateURL
		metric := dto.Metrics{ID: name, MType: mtyp}

		switch n.GetMetric().MetricType() {
		case types.COUNTER:
			i64, err := strconv.ParseInt(*value, 10, 64)
			if err != nil {
				panic(err)
			}
			metric.Delta = &i64
		case types.GAUGE:
			f64, err := strconv.ParseFloat(*value, 64)
			if err != nil {
				panic(err)
			}
			metric.Value = &f64
		}

		var b bytes.Buffer
		if _, err := easyjson.MarshalToWriter(metric, &b); err != nil {
			panic(err)
		}

		request, err := http.NewRequest(http.MethodPost, path, &b)
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
				logger.Log.Debug(
					"func postDo",
					zap.String("error", fmt.Sprintf("%v", p)),
				)
			}
		}
	}()

	request.Header.Add("Content-Type", "application/json")
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
