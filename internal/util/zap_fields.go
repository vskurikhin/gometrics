/*
 * This file was last modified at 2024-03-01 21:40 by Victor N. Skurikhin.
 * zap_fields.go
 * $Id$
 */

package util

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/dto"
	"go.uber.org/zap"
)

type zapFields []zap.Field

//goland:noinspection GoExportedFuncWithUnexportedType
func MakeZapFields() zapFields {
	return make(zapFields, 0)
}

func (zf *zapFields) Append(key string, metric interface{}) {
	*zf = append(*zf, zap.String(key, fmt.Sprintf("%+v", metric)))
}

func (zf *zapFields) AppendInt(key string, metric interface{}) {
	*zf = append(*zf, zap.String(key, fmt.Sprintf("%d", metric)))
}

func (zf *zapFields) AppendFloat(key string, metric interface{}) {
	*zf = append(*zf, zap.String(key, fmt.Sprintf("%.12f", metric)))
}

//goland:noinspection GoExportedFuncWithUnexportedType
func ZapFieldsMetric(metric *dto.Metrics) zapFields {
	zapFields := MakeZapFields()
	zapFields.Append("metric", metric)

	if metric.Value != nil {
		zapFields.AppendFloat("Value", *metric.Value)
	}
	if metric.Delta != nil {
		zapFields.AppendInt("Delta", *metric.Delta)
	}
	return zapFields
}
