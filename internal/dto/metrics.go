/*
 * This file was last modified at 2024-02-28 23:19 by Victor N. Skurikhin.
 * metrics.go
 * $Id$
 */

package dto

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
