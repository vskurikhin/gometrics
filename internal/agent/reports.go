/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * reports.go
 * $Id$
 */

package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/vskurikhin/gometrics/internal/crypto"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
)

func Reports(cfg env.Config, enabled []types.Name) {

	client := http.Client{}
	for {
		reports(cfg, enabled, &client)
	}
}

func reports(cfg env.Config, enabled []types.Name, client *http.Client) {

	time.Sleep(cfg.ReportInterval() * time.Second)

	metrics := make(dto.Metrics, 0)

	for _, i := range enabled {
		metric := getMetric(i)
		if metric != nil {
			metrics = append(metrics, *metric)
		}
	}
	request, err := NewRequest(cfg, metrics)

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
func NewRequest(cfg env.Config, metrics dto.Metrics) (*http.Request, error) {

	marshaledBuffer := &bytes.Buffer{}

	if _, err := easyjson.MarshalToWriter(metrics, marshaledBuffer); err != nil {
		return nil, err
	}

	compressBuffer := &bytes.Buffer{}

	gz, err := gzip.NewWriterLevel(compressBuffer, gzip.BestCompression)

	if err != nil {
		//nolint:multichecker,errcheck
		_, _ = io.WriteString(marshaledBuffer, err.Error())
		return nil, err
	}
	//nolint:multichecker,errcheck
	_, _ = gz.Write(marshaledBuffer.Bytes())
	//nolint:multichecker,errcheck
	_ = gz.Close()
	crypt := crypto.GetAgentCrypto(cfg)
	buffer := &bytes.Buffer{}

	// TODO переработать поблочную обработку
	if buf, err := crypt.EncryptRSA(compressBuffer.Bytes()); err != nil {
		logger.Log.Debug("encrypt fail", zap.String("error", fmt.Sprintf("%v", err)))
		buffer = bytes.NewBuffer(compressBuffer.Bytes())
	} else {
		buffer = bytes.NewBuffer(buf)
	}

	path := *cfg.URLHost() + env.UpdatesURL
	request, err := http.NewRequest(http.MethodPost, path, buffer)

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
		if err != nil {
			logger.Log.Debug("post fail", zap.String("error", fmt.Sprintf("%v", err)))
		}
		if response != nil {
			//nolint:multichecker,errcheck
			_ = response.Body.Close()
		}
	}()
	return err
}
