/*
 * This file was last modified at 2024-02-04 13:11 by Victor N. Skurikhin.
 * mem_storage.go
 * $Id$
 */

package memory

import (
	"github.com/vskurikhin/gometrics/api/names"
	"sync"
)

type MemStorage struct {
	mu      sync.Mutex
	Metrics map[string]interface{}
}

var mem = new(MemStorage)

func init() {
	mem.Metrics = make(map[string]interface{})
}

func Instance() *MemStorage {
	return mem
}

func (m *MemStorage) Get(name names.Names) (interface{}, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	return m.Metrics[name.String()], nil
}

func (m *MemStorage) Put(name names.Names, value interface{}) error {

	m.mu.Lock()
	defer m.mu.Unlock()
	m.Metrics[name.String()] = value

	return nil
}
