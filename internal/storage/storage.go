/*
 * This file was last modified at 2024-03-02 14:18 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

type Storage interface {
	Get(name string) *string

	Put(name string, value *string)

	ReadFromFile(fileName string)

	SaveToFile(fileName string)
}
