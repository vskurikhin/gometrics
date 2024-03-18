/*
 * This file was last modified at 2024-03-18 18:57 by Victor N. Skurikhin.
 * mem_storage.go
 * $Id$
 */

package memory

import (
	"encoding/json"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/util"
	"io"
	"os"
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

func Instance() storage.Storage {
	return mem
}

func (m *MemStorage) Get(name string) *string {

	m.RLock()
	defer m.RUnlock()

	return m.metrics[name]
}

func (m *MemStorage) GetCounter(name string) *string {
	return m.Get(name)
}

func (m *MemStorage) GetGauge(name string) *string {
	return m.Get(name)
}

func (m *MemStorage) Put(name string, value *string) {

	m.Lock()
	defer m.Unlock()
	m.metrics[name] = value
}

func (m *MemStorage) PutCounter(name string, value *string) {
	m.Put(name, value)
}

func (m *MemStorage) PutGauge(name string, value *string) {
	m.Put(name, value)
}

func (m *MemStorage) ReadFromFile(fileName string) {

	zapFields := util.MakeZapFields()
	zapFields.AppendString("fileName", &fileName)

	n, err := m.readFromFile(zapFields, fileName)
	if err != nil {
		panic(err)
	}
	zapFields.AppendInt("read bytes", n)
	zapFields.Append("m.metrics", m.metrics)
	logger.Log.Debug("in ReadFromFile", zapFields.Slice()...)
}

func (m *MemStorage) readFromFile(zf util.ZapFields, fileName string) (int, error) {

	m.Lock()
	defer m.Unlock()

	file, err := os.Open(fileName)

	if err != nil {
		zf.Append("error", err)
		logger.Log.Error("in ReadFromFile", zf.Slice()...)
		return 0, nil
	}
	defer file.Close()

	buf, err := io.ReadAll(file)

	if err != nil {
		zf.Append("error", err)
		logger.Log.Error("in ReadFromFile", zf.Slice()...)
		return 0, nil
	}
	err = json.Unmarshal(buf, &m.metrics)

	if err != nil {
		zf.Append("error", err)
		logger.Log.Error("in ReadFromFile", zf.Slice()...)
		return 0, nil
	}
	return len(buf), nil
}

func (m *MemStorage) SaveToFile(fileName string) {

	zapFields := util.MakeZapFields()
	zapFields.AppendString("fileName", &fileName)

	out, n := m.saveToFile(zapFields, fileName)

	jsonMetrics := string(out)
	zapFields.AppendInt("write bytes", n)
	zapFields.AppendString("m.metrics", &jsonMetrics)
	logger.Log.Debug("in SaveToFile", zapFields.Slice()...)
}

func (m *MemStorage) saveToFile(zf util.ZapFields, fileName string) ([]byte, int) {

	m.Lock()
	defer m.Unlock()
	out, err := json.Marshal(m.metrics)

	if err != nil {
		zf.Append("error", err)
		logger.Log.Error("in SaveToFile", zf.Slice()...)
		return nil, 0
	}
	var flag int
	if env.Server.StoreInterval() == 0 {
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC | os.O_SYNC
	} else {
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	}
	file, err := os.OpenFile(fileName, flag, 0640)

	if err != nil {
		zf.Append("error", err)
		logger.Log.Error("in SaveToFile", zf.Slice()...)
		return nil, 0
	}
	defer file.Close()

	n, err := file.Write(out)

	if err != nil {
		zf.Append("error", err)
		logger.Log.Error("in SaveToFile", zf.Slice()...)
		return nil, 0
	}
	return out, n
}
