/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * reports.go
 * $Id$
 */

package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vskurikhin/gometrics/internal/crypto"
	"github.com/vskurikhin/gometrics/internal/util"

	"github.com/mailru/easyjson"
	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
)

func Reports(ctx context.Context, cfg env.Config, enabled []types.Name) {

	client := http.Client{}
	time.Sleep(cfg.PollInterval())
	grpcReports(cfg, enabled)
	for {
		select {
		case <-ctx.Done():
			grpcReports(cfg, enabled)
			return
		default:
			go grpcReports(cfg, enabled)
			go reports(cfg, enabled, &client)
			time.Sleep(cfg.ReportInterval())
		}
	}
}

func reports(cfg env.Config, enabled []types.Name, client *http.Client) {

	metrics := make(dto.Metrics, 0)

	for _, i := range enabled {
		metric := getMetric(i)
		if metric != nil {
			metrics = append(metrics, *metric)
		}
	}
	request, err := newRequest(cfg, metrics)
	util.IfErrorThenPanic(err)
	err = postDo(client, request)

	for i := 0; err != nil && isUpperBound(i, cfg.ReportInterval()); i++ {
		time.Sleep(time.Duration(1<<i) * time.Second)
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
			util.IfErrorThenPanic(err)
			metric.Delta = &i64
		case types.GAUGE:
			f64, err := strconv.ParseFloat(*value, 64)
			util.IfErrorThenPanic(err)
			metric.Value = &f64
		}
		return &metric
	}
	return nil
}

//goland:noinspection GoUnhandledErrorResult
func newRequest(cfg env.Config, metrics dto.Metrics) (*http.Request, error) {

	var marshalBuffer, requestBuffer bytes.Buffer

	if _, err := easyjson.MarshalToWriter(metrics, &marshalBuffer); err != nil {
		return nil, err
	}

	crypt := crypto.GetAgentCrypto()
	gz, err := gzip.NewWriterLevel(&requestBuffer, gzip.BestSpeed)

	if err != nil {
		return nil, err
	}
	body := marshalBuffer.Bytes()
	secretKey, be, err := crypt.EncryptAES(body)

	if err != nil {
		//nolint:multichecker,errcheck
		_, _ = gz.Write(body)
	} else {
		//nolint:multichecker,errcheck
		_, _ = gz.Write(be)
	}
	//nolint:multichecker,errcheck
	_ = gz.Close()

	path := *cfg.URLHost() + env.UpdatesURL
	request, err := http.NewRequest(http.MethodPost, path, &requestBuffer)

	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Encoding", "gzip")
	key := encryptRSAAndBase64StdEncoding(crypt, secretKey)
	request.Header.Add("X-Content-Encrypting", key)
	logger.Log.Debug("post EncryptRSA", zap.String("key", key))
	ip := cfg.Property().OutboundIP()
	request.Header.Add("X-Real-IP", ip.String())

	return request, nil
}

func postDo(client *http.Client, request *http.Request) error {

	response, err := client.Do(request)

	defer func() {
		if response != nil {
			//nolint:multichecker,errcheck
			_ = response.Body.Close()
		}
	}()
	return err
}

func encryptRSAAndBase64StdEncoding(crypt crypto.Crypto, plain []byte) string {
	if bs, err := crypt.EncryptRSA(plain); err != nil {
		logger.Log.Debug("post EncryptRSA",
			zap.String("error", fmt.Sprintf("%v", err)),
		)
		return ""
	} else {
		return base64.StdEncoding.EncodeToString(bs)
	}
}

func isUpperBound(index int, duration time.Duration) bool {
	result := time.Duration((index*(index+1)*(2*index+1))/6) * time.Second
	return result < duration
}
