/*
 * This file was last modified at 2024-02-29 23:38 by Victor N. Skurikhin.
 * zap_fields.go
 * $Id$
 */

package util

import (
	"fmt"
	"go.uber.org/zap"
)

type zapFields []zap.Field

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
