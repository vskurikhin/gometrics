/*
 * This file was last modified at 2024-03-18 18:57 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

type Storage interface {
	Get(name string) *string

	GetCounter(name string) *string

	GetGauge(name string) *string

	Put(name string, value *string)

	PutCounter(name string, value *string)

	PutGauge(name string, value *string)

	ReadFromFile(fileName string)

	SaveToFile(fileName string)
}
