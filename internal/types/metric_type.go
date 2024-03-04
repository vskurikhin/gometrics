/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * metric_type.go
 * $Id$
 */

package types

import (
	"strconv"
	"strings"
)

type (
	Counter int64
	Gauge   float64
)

type Types int

const (
	COUNTER Types = iota
	GAUGE
)

var lower []string
var types = [...]string{"Counter", "Gauge"}

func (t Types) String() string {
	return types[t]
}

func (t Types) URLPath() string {
	return lower[t]
}

func (t Types) Eq(s string) bool {
	return strings.EqualFold(strings.ToLower(s), strings.ToLower(t.String()))
}

func (t Types) ParseValue(s string) (interface{}, error) {

	switch t {
	case COUNTER:
		return strconv.Atoi(s)
	case GAUGE:
		return strconv.ParseFloat(s, 64)
	default:
		return strconv.ParseFloat(s, 64)
	}
}
