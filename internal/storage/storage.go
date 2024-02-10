/*
 * This file was last modified at 2024-02-11 00:04 by Victor N. Skurikhin.
 * storage.go
 * $Id$
 */

package storage

import (
	"github.com/vskurikhin/gometrics/internal/types"
)

type Storage interface {
	Get(name types.Name) (interface{}, error)

	Put(name types.Name, value interface{}) error
}
