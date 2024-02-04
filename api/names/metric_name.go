/*
 * This file was last modified at 2024-02-04 13:06 by Victor N. Skurikhin.
 * metric_name.go
 * $Id$
 */

package names

import (
	"slices"
	"strings"
)

type Names int

const (
	_none Names = iota
	Alloc
	BuckHashSys
	Frees
)

var l []string
var s = [...]string{"Alloc", "BuckHashSys", "Frees"}

func init() {
	for i := range s {
		l = append(l, strings.ToLower(s[i]))
	}
}

func Lookup(s string) Names {
	i := lookup(s)
	if i < 0 {
		return 1
	}
	return Names(i)
}

func lookup(s string) int {
	return slices.IndexFunc(l, func(e string) bool {
		return e == strings.ToLower(s)
	})
}

func (t Names) String() string {
	return s[t]
}

func (t Names) URLPath() string {
	return l[t]
}

func (t Names) Eq(s string) bool {
	return strings.EqualFold(strings.ToLower(s), strings.ToLower(t.String()))
}
