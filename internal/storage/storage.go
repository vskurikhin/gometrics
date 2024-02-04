/*
 * This file was last modified at 2024-02-04 15:21 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

import "github.com/vskurikhin/gometrics/api/names"

type Storage interface {
	Get(name names.Names) (interface{}, error)

	Put(name names.Names, value interface{}) error
}
