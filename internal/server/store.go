/*
 * This file was last modified at 2024-03-02 14:24 by Victor N. Skurikhin.
 * store.go
 * $Id$
 */

package server

import "github.com/vskurikhin/gometrics/internal/storage/memory"

var store = memory.Instance()
