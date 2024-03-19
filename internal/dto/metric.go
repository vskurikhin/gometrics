/*
 * This file was last modified at 2024-03-19 09:31 by Victor N. Skurikhin.
 * metric.go
 * $Id$
 */

package dto

import "strconv"

type Metric struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m *Metric) CalcDelta(value *string) int64 {

	var err error
	var i64 int64

	if value != nil {
		i64, err = strconv.ParseInt(*value, 10, 64)
	}
	if err != nil {
		panic(err)
	}
	if m.Delta != nil {
		i64 += *m.Delta
	}
	return i64
}
