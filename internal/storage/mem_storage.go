/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * mem_storage.go
 * $Id$
 */

package storage

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/types"
	"github.com/vskurikhin/gometrics/internal/util"
)

type MemStorage struct {
	sync.RWMutex
	Metrics       map[string]*string
	StoreInterval int
}

func (m *MemStorage) Get(name string) *string {
	return m.get(name)
}

func (m *MemStorage) GetCounter(name string) *string {
	return m.get(name)
}

func (m *MemStorage) GetGauge(name string) *string {
	return m.get(name)
}

func (m *MemStorage) get(name string) *string {

	m.RLock()
	defer m.RUnlock()

	return m.Metrics[name]
}

// Deprecated: Put is deprecated.
func (m *MemStorage) Put(name string, value *string) {
	m.put(name, value)
}

func (m *MemStorage) PutCounter(name string, value *string) {
	m.put(name, value)
}

func (m *MemStorage) PutGauge(name string, value *string) {
	m.put(name, value)
}

func (m *MemStorage) put(name string, value *string) {

	m.Lock()
	defer m.Unlock()
	m.Metrics[name] = value
}

func (m *MemStorage) PutSlice(metrics dto.Metrics) {

	for _, metric := range metrics {

		num := types.Lookup(metric.ID)
		var name string

		if num > 0 {
			name = num.String()
		} else {
			name = metric.ID
		}
		value := m.Get(name)

		switch {
		case types.GAUGE.Eq(metric.MType):
			value := fmt.Sprintf("%.12f", *metric.Value)
			m.PutGauge(name, &value)
		case types.COUNTER.Eq(metric.MType):
			*metric.Delta = metric.CalcDelta(value)
			value := fmt.Sprintf("%d", *metric.Delta)
			m.PutCounter(name, &value)
		}
	}
}

func (m *MemStorage) ReadFromFile(fileName string) {

	zapFields := util.MakeZapFields()
	zapFields.AppendString("fileName", &fileName)

	n, err := m.readFromFile(zapFields, fileName)
	if err != nil {
		panic(err)
	}
	zapFields.AppendInt("read bytes", n)
	zapFields.Append("m.Metrics", m.Metrics)
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
	err = json.Unmarshal(buf, &m.Metrics)

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
	zapFields.AppendString("m.Metrics", &jsonMetrics)
	logger.Log.Debug("in SaveToFile", zapFields.Slice()...)
}

func (m *MemStorage) saveToFile(zf util.ZapFields, fileName string) ([]byte, int) {

	m.Lock()
	defer m.Unlock()
	out, err := json.Marshal(m.Metrics)

	if err != nil {
		zf.Append("error", err)
		logger.Log.Error("in SaveToFile", zf.Slice()...)
		return nil, 0
	}
	var flag int
	if m.StoreInterval == 0 {
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC | os.O_SYNC
	} else {
		flag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	}
	file, err := os.OpenFile(fileName, flag, 0640)
	for i := 1; err != nil && i < 6; i += 2 {
		time.Sleep(time.Duration(i) * time.Second)
		logger.Log.Debug("retry open file",
			zap.String("error", fmt.Sprintf("%v", err)),
			zap.String("time", fmt.Sprintf("%v", time.Now())),
		)
		file, err = os.OpenFile(fileName, flag, 0640)
	}

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
