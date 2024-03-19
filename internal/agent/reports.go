/*
 * This file was last modified at 2024-03-19 12:04 by Victor N. Skurikhin.
 * reports.go
 * $Id$
 */

package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Reports(enabled []types.Name) {

	client := http.Client{}
	for {
		reports(enabled, &client)
	}
}

func reports(enabled []types.Name, client *http.Client) {

	time.Sleep(env.Agent.ReportInterval() * time.Second)

	metrics := make(dto.Metrics, 0)

	for _, i := range enabled {
		metric := getMetric(i)
		if metric != nil {
			metrics = append(metrics, *metric)
		}
	}
	request, err := newRequest(metrics)

	if err != nil {
		panic(err)
	}

	err = postDo(client, request)

	for i := 1; err != nil && i < 6; i += 2 {
		time.Sleep(time.Duration(i) * time.Second)
		logger.Log.Debug("retry post",
			zap.String("error", fmt.Sprintf("%v", err)),
			zap.String("time", fmt.Sprintf("%v", time.Now())),
		)
		err = postDo(client, request)
	}
}

func getMetric(n types.Name) *dto.Metric {

	name := n.GetMetric().String()
	var value *string

	switch n.GetMetric().MetricType() {
	case types.COUNTER:
		value = store.GetCounter(name)
	case types.GAUGE:
		value = store.GetGauge(name)
	default:
		value = store.Get(name)
	}

	if value != nil {

		mtyp := n.
			GetMetric().
			MetricType().
			URLPath()
		metric := dto.Metric{ID: name, MType: mtyp}

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

		return &metric
	}
	return nil
}

//goland:noinspection GoUnhandledErrorResult
func newRequest(metrics dto.Metrics) (*http.Request, error) {

	var b1, b2 bytes.Buffer

	if _, err := easyjson.MarshalToWriter(metrics, &b1); err != nil {
		return nil, err
	}
	gz, err := gzip.NewWriterLevel(&b2, gzip.BestSpeed)

	if err != nil {
		io.WriteString(&b1, err.Error())
		return nil, err
	}
	gz.Write(b1.Bytes())
	gz.Close()

	path := *env.Agent.URLHost() + names.UpdatesURL
	request, err := http.NewRequest(http.MethodPost, path, &b2)

	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Encoding", "gzip")

	return request, nil
}

func postDo(client *http.Client, request *http.Request) error {

	response, err := client.Do(request)

	defer func() {
		if response != nil {
			//goland:noinspection GoUnhandledErrorResult
			response.Body.Close()
		}
	}()
	return err
}
