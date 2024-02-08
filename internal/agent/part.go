/*
 * This file was last modified at 2024-02-08 08:57 by Victor N. Skurikhin.
 * part.go
 * $Id$
 */

package agent

type part string

type urlPart interface {
	URLPath() string
}

func (s part) URLPath() string {
	return string(s)
}
