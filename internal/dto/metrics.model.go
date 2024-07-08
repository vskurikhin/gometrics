/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * metrics.model.go
 * $Id$
 */

package dto

import "bytes"

// Metrics метрики
// easyjson:json
type Metrics []Metric

func (m *Metrics) String() string {
	if m != nil {
		var buffer bytes.Buffer
		buffer.WriteRune('[')
		for i, metric := range *m {
			if i != 0 {
				buffer.WriteString(", ")
			}
			buffer.WriteString(metric.String())
		}
		buffer.WriteRune(']')
		return buffer.String()
	}
	return "<nil>"
}
