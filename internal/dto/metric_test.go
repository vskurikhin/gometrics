/*
 * This file was last modified at 2024-06-10 11:18 by Victor N. Skurikhin.
 * metric_test.go
 * $Id$
 */

package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
	"os"
	"testing"
)

//goland:noinspection GoUnhandledErrorResult
func TestMetric(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]interface{}
	}{
		{
			name: "positive test #1",
			input: map[string]interface{}{
				"id":    "Alloc",
				"type":  "gauge",
				"value": 1.1,
			},
		},
		{
			name: "positive test #2",
			input: map[string]interface{}{
				"id":    "Alloc",
				"type":  "gauge",
				"value": "1.1",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			metric := new(Metric)

			body, _ := json.Marshal(test.input)
			metric.UnmarshalJSON(body)
			metric.MarshalJSON()
			w := new(jwriter.Writer)
			metric.MarshalEasyJSON(w)
			easyjson.UnmarshalFromReader(bytes.NewReader(body), metric)
		})
	}
}

func TestMetricCalcDeltaPositive(t *testing.T) {
	var d int64
	p := "1"
	tests := []struct {
		name  string
		input *string
	}{
		{
			name:  "positive test #1",
			input: &p,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			metric := Metric{Delta: &d}

			res := metric.CalcDelta(test.input)
			fmt.Fprintf(os.Stderr, "res: %v\n", res)
		})
	}
}

func TestMetricCalcDeltaNegative(t *testing.T) {
	n1 := ""
	n2 := "1.1"
	tests := []struct {
		name  string
		input *string
	}{
		{
			name:  "negative test #1",
			input: &n1,
		},
		{
			name:  "negative test #2",
			input: &n2,
		},
		{
			name:  "negative test #3",
			input: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			defer func() {
				if p := recover(); p != nil {
					fmt.Fprintf(os.Stderr, "recover: %v\n", p)
				}
			}()

			var metric Metric

			res := metric.CalcDelta(test.input)
			fmt.Fprintf(os.Stderr, "res: %v\n", res)
		})
	}
}
