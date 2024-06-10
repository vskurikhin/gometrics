/*
 * This file was last modified at 2024-06-10 13:54 by Victor N. Skurikhin.
 * str.go
 * $Id$
 */

package util

func Str(s *string) string {
	if s != nil {
		return *s
	}
	return "<nil>"
}
