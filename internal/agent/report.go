/*
 * This file was last modified at 2024-04-05 09:09 by Victor N. Skurikhin.
 * report.go
 * $Id$
 */

package agent

import (
	"bytes"
	"compress/gzip"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/types"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Deprecated: Report is deprecated.
func Report(enabled []types.Name) {

	client := http.Client{}
	for {
		report(enabled, &client)
	}
}

// Deprecated: report is deprecated.
func report(enabled []types.Name, client *http.Client) {

	time.Sleep(env.Agent.ReportInterval() * time.Second)

	for _, i := range enabled {
		post(i, client)
	}
}

// Deprecated: post is deprecated.
func post(n types.Name, client *http.Client) {

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
		path := *env.Agent.URLHost() + env.UpdateURL
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
		var b1 bytes.Buffer

		if _, err := easyjson.MarshalToWriter(metric, &b1); err != nil {
			panic(err)
		}
		var b2 bytes.Buffer
		gz, err := gzip.NewWriterLevel(&b2, gzip.BestSpeed)

		if err != nil {
			io.WriteString(&b1, err.Error())
			return
		}
		gz.Write(b1.Bytes())
		gz.Close()
		request, err := http.NewRequest(http.MethodPost, path, &b2)

		if err != nil {
			panic(err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Content-Encoding", "gzip")
		postDo(client, request)
	}
}
