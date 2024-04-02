/*
 * This file was last modified at 2024-03-19 10:29 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

import "github.com/vskurikhin/gometrics/internal/dto"

type Storage interface {
	Get(name string) *string

	GetCounter(name string) *string

	GetGauge(name string) *string

	Put(name string, value *string)

	PutCounter(name string, value *string)

	PutGauge(name string, value *string)

	PutSlice(metrics dto.Metrics)

	ReadFromFile(fileName string)

	SaveToFile(fileName string)
}
