/*
 * This file was last modified at 2024-06-10 09:35 by Victor N. Skurikhin.
 * zap_fields_test.go
 * $Id$
 */

package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gometrics/internal/dto"
	"testing"
)

func TestMakeZapFields(t *testing.T) {
	var f = 0.1
	var i int64 = 1
	var key = "key"
	var tests = []struct {
		name string
		do   func() interface{}
	}{
		{name: "positive test #0", do: func() interface{} { return MakeZapFields() }},
		{name: "positive test #1", do: func() interface{} { z := MakeZapFields(); z.Append(key, nil); return z }},
		{name: "positive test #2", do: func() interface{} { z := MakeZapFields(); z.AppendInt(key, i); return z }},
		{name: "positive test #3", do: func() interface{} { z := MakeZapFields(); z.AppendFloat(key, f); return z }},
		{name: "positive test #4", do: func() interface{} { z := MakeZapFields(); z.AppendString(key, &key); return z }},
		{name: "positive test #5", do: func() interface{} { z := MakeZapFields(); z.AppendString(key, nil); return z }},
		{name: "positive test #6", do: func() interface{} { z := MakeZapFields(); return z.Slice() }},
		{name: "positive test #7", do: func() interface{} { return ZapFieldsMetric(&dto.Metric{Delta: &i}) }},
		{name: "positive test #8", do: func() interface{} { return ZapFieldsMetric(&dto.Metric{Value: &f}) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.do()
			assert.NotNil(t, got)
		})
	}
}
