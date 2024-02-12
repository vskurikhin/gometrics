/*
 * This file was last modified at 2024-02-11 13:51 by Victor N. Skurikhin.
 * mem_storage.go
 * $Id$
 */

package memory

import (
	"sync"
)

type MemStorage struct {
	sync.RWMutex
	metrics map[string]*string
}

var mem = new(MemStorage)

func init() {
	mem.metrics = make(map[string]*string)
}

func Instance() *MemStorage {
	return mem
}

func (m *MemStorage) Get(name string) *string {

	m.RLock()
	defer m.RUnlock()

	return m.metrics[name]
}

func (m *MemStorage) Put(name string, value *string) {

	m.Lock()
	defer m.Unlock()
	m.metrics[name] = value
}
