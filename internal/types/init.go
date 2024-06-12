/*
 * This file was last modified at 2024-02-25 12:11 by Victor N. Skurikhin.
 * init.go
 * $Id$
 */

// Package types типы и имена метрик
package types

import "strings"

func init() {

	for i := range types {
		lower = append(lower, strings.ToLower(types[i]))
	}
}
