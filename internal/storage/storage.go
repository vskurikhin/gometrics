/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

type Storage interface {
	Get(name string) *string

	Put(name string, value *string)
}
