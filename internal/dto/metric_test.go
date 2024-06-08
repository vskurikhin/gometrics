/*
 * This file was last modified at 2024-05-28 21:58 by Victor N. Skurikhin.
 * metric_test.go
 * $Id$
 */

package dto

import (
	"bytes"
	"encoding/json"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jwriter"
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
