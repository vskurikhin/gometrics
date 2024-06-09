/*
 * This file was last modified at 2024-03-19 11:03 by Victor N. Skurikhin.
 * metric.go
 * $Id$
 */

// Package dto Data Transfer Object
package dto

import "strconv"

// Metric метрика
type Metric struct {
	ID    string   `json:"id"`              // ID имя метрики
	MType string   `json:"type"`            // MType параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Delta значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Value значение метрики в случае передачи gauge
}

// CalcDelta расчёт дельты
func (m *Metric) CalcDelta(value *string) int64 {

	var i64 int64

	if value != nil {
		i64 = parseInt(value)
	}
	if m.Delta != nil {
		i64 += *m.Delta
	}
	return i64
}

func parseInt(value *string) int64 {
	if i64, err := strconv.ParseInt(*value, 10, 64); err != nil {
		panic(err)
	} else {
		return i64
	}
}
