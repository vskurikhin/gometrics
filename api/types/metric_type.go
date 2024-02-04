/*
 * This file was last modified at 2024-02-04 13:02 by Victor N. Skurikhin.
 * metric_type.go
 * $Id$
 */

package types

import (
	"errors"
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

var l []string
var s = [...]string{"Counter", "Gauge"}

func init() {
	for i := range s {
		l = append(l, strings.ToLower(s[i]))
	}
}

func (t Types) String() string {
	return s[t]
}

func (t Types) URLPath() string {
	return l[t]
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
	}
	return nil, errors.New("can't parse " + s)
}
