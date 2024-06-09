/*
 * This file was last modified at 2024-06-10 14:28 by Victor N. Skurikhin.
 * if_error_then_panic.go
 * $Id$
 */

package util

func IfErrorThenPanic(e error) {
	if e != nil {
		panic(e)
	}
}
