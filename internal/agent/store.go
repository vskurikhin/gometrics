/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * store.go
 * $Id$
 */

package agent

import "github.com/vskurikhin/gometrics/internal/storage/memory"

var store = memory.Instance()
