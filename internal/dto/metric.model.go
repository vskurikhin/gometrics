/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * metric.model.go
 * $Id$
 */

// Package dto Data Transfer Object
package dto

import (
	"fmt"
	"strconv"

	pb "github.com/vskurikhin/gometrics/proto"
)

// Metric метрика
type Metric struct {
	ID    string   `json:"id"`              // ID имя метрики
	MType string   `json:"type"`            // MType параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Delta значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Value значение метрики в случае передачи gauge
}

func FromRequest(protoBuf *pb.Metric) *Metric {
	delta := protoBuf.GetDelta()
	value := protoBuf.GetValue()
	return &Metric{
		ID:    protoBuf.GetId(),
		MType: protoBuf.GetType(),
		Delta: &delta,
		Value: &value,
	}
}

func FromValueRequest(protoBuf *pb.MetricRequestValue) *Metric {
	return &Metric{
		ID:    protoBuf.GetId(),
		MType: protoBuf.GetType(),
	}
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

func (m *Metric) ToResponse() *pb.Metric {

	result := new(pb.Metric)
	result.Id = m.ID
	result.Type = m.MType

	if m.Delta != nil {
		result.Delta = *m.Delta
	}
	if m.Value != nil {
		result.Value = *m.Value
	}
	return result
}

func (m *Metric) String() string {
	if m != nil {
		var delta, value = "<nil>", "<nil>"
		if m.Delta != nil {
			delta = fmt.Sprintf("%d", *m.Delta)
		}
		if m.Value != nil {
			value = fmt.Sprintf("%e", *m.Value)
		}
		return fmt.Sprintf("{ID:%s MType:%s Delta:%s Value:%s}", m.ID, m.MType, delta, value)
	}
	return "<nil>"
}

func parseInt(value *string) int64 {
	if i64, err := strconv.ParseInt(*value, 10, 64); err != nil {
		panic(err)
	} else {
		return i64
	}
}
