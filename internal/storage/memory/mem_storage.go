/*
 * This file was last modified at 2024-02-08 21:41 by Victor N. Skurikhin.
 * mem_storage.go
 * $Id$
 */

package memory

import (
	"sync"
)

type MemStorage struct {
	mu      sync.RWMutex
	Metrics map[string]*string
}

var mem = new(MemStorage)

func init() {
	mem.Metrics = make(map[string]*string)
}

func Instance() *MemStorage {
	return mem
}

func (m *MemStorage) Get(name string) *string {

	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.Metrics[name]
}

func (m *MemStorage) Put(name string, value *string) {

	m.mu.Lock()
	defer m.mu.Unlock()
	m.Metrics[name] = value
}
