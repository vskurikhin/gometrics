/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * zap_fields.go
 * $Id$
 */

package util

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/dto"
)

type zapFields []zap.Field

type ZapFields interface {
	Append(key string, metric interface{})

	AppendInt(key string, metric interface{})

	AppendFloat(key string, metric interface{})

	AppendString(key string, s *string)

	Slice() []zap.Field
}

//goland:noinspection GoExportedFuncWithUnexportedType
func MakeZapFields() ZapFields {
	result := make(zapFields, 0)
	return &result
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

func (zf *zapFields) AppendString(key string, s *string) {
	if s != nil {
		*zf = append(*zf, zap.String(key, *s))
	} else {
		*zf = append(*zf, zap.String(key, "nil"))
	}
}

func (zf *zapFields) Slice() []zap.Field {
	result := make([]zap.Field, len(*zf))
	copy(result, *zf)
	return result
}

//goland:noinspection GoExportedFuncWithUnexportedType
func ZapFieldsMetric(metric *dto.Metric) ZapFields {
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
