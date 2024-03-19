/*
 * This file was last modified at 2024-03-19 10:31 by Victor N. Skurikhin.
 * pgs_storage.go
 * $Id$
 */

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gometrics/internal/dto"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/types"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type PgsStorage struct {
	memory storage.Storage
	pool   *pgxpool.Pool
}

const (
	sqlInsertCounter = `INSERT INTO metric (name, type, counter)
				VALUES ($1, 'counter', $2)
				ON CONFLICT(name) 
				DO UPDATE SET
				counter = $2`
	sqlInsertGauge = `INSERT INTO metric (name, type, gauge)
				VALUES ($1, 'gauge', $2)
				ON CONFLICT(name) 
				DO UPDATE SET
				gauge = $2`
)

func New(memory storage.Storage, pool *pgxpool.Pool) *PgsStorage {
	return &PgsStorage{
		memory: memory,
		pool:   pool,
	}
}

// Deprecated: Get is deprecated.
func (p *PgsStorage) Get(name string) *string {
	return p.memory.Get(name)
}

func (p *PgsStorage) GetCounter(name string) *string {

	// cache
	value := p.memory.Get(name)

	if value != nil {
		return value
	}
	row, err := p.getSQL("SELECT counter FROM metric WHERE name = $1 AND type = 'counter'", name)

	if err != nil {
		panic(err)
	}
	var counter sql.NullInt64

	err = row.Scan(&counter)

	if err != nil {
		return nil
	}

	if counter.Valid {
		result := fmt.Sprintf("%d", counter.Int64)
		return &result
	}
	return nil
}

func (p *PgsStorage) GetGauge(name string) *string {

	// cache
	value := p.memory.GetGauge(name)

	if value != nil {
		return value
	}
	row, err := p.getSQL("SELECT gauge FROM metric WHERE name = $1 AND type = 'gauge'", name)

	if err != nil {
		panic(err)
	}
	var gauge sql.NullFloat64

	err = row.Scan(&gauge)

	if err != nil {
		return nil
	}
	if gauge.Valid {
		result := fmt.Sprintf("%.3f", gauge.Float64)
		if len(result) > 0 && result[len(result)-1] == '0' {
			result = result[:len(result)-1]
		}
		return &result
	}
	return nil
}

func (p *PgsStorage) getSQL(sql, name string) (pgx.Row, error) {

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func getSQL", zap.String("error", fmt.Sprintf("%v", p)))
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := p.pool.Acquire(ctx)
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	if conn == nil || err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	return conn.QueryRow(ctx, sql, name), nil
}

// Deprecated: Put is deprecated.
func (p *PgsStorage) Put(name string, value *string) {
	p.memory.Put(name, value)
}

func (p *PgsStorage) PutCounter(name string, value *string) {

	counter, err := strconv.Atoi(*value)

	if err != nil {
		panic(err)
	}
	err = p.putSQL(sqlInsertCounter, name, counter)
	if err != nil {
		panic(err)
	}
	p.memory.Put(name, value)
}

func (p *PgsStorage) PutGauge(name string, value *string) {

	gauge, err := strconv.ParseFloat(*value, 64)

	if err != nil {
		panic(err)
	}
	err = p.putSQL(sqlInsertGauge, name, gauge)
	if err != nil {
		panic(err)
	}
	p.memory.Put(name, value)
}

func (p *PgsStorage) putSQL(sql, name string, value interface{}) error {

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func putSQL", zap.String("error", fmt.Sprintf("%v", p)))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := p.pool.Acquire(ctx)
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	if conn == nil || err != nil {
		return fmt.Errorf("%v", err)
	}

	_, err = conn.Exec(ctx, "INSERT INTO name (name) VALUES ($1) ON CONFLICT DO NOTHING", name)

	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, sql, name, value)

	if err != nil {
		return err
	}
	return nil
}

func (p *PgsStorage) PutSlice(metrics dto.Metrics) {

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func PutSlice", zap.String("error", fmt.Sprintf("%v", p)))
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := p.pool.Acquire(ctx)
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()
	if conn == nil || err != nil {
		panic(err)
	}

	tx, err := conn.Begin(ctx)

	if err != nil {
		panic(err)
	}
	for _, metric := range metrics {

		num := types.Lookup(metric.ID)
		var name string

		if num > 0 {
			name = num.String()
		} else {
			name = metric.ID
		}
		value := p.Get(name)

		var sqlCommand string
		switch {
		case types.GAUGE.Eq(metric.MType):
			sqlCommand = sqlInsertGauge
			v := fmt.Sprintf("%.12f", *metric.Value)
			value = &v
		case types.COUNTER.Eq(metric.MType):
			sqlCommand = sqlInsertCounter
			*metric.Delta = metric.CalcDelta(value)
			v := fmt.Sprintf("%d", *metric.Delta)
			value = &v
		}
		_, err = tx.Exec(ctx, sqlCommand, name, value)
	}
	if err != nil {
		panic(err)
	}
	err = tx.Commit(ctx)

	if err != nil {
		panic(err)
	}
	p.memory.PutSlice(metrics)
}

func (p *PgsStorage) ReadFromFile(fileName string) {
	p.memory.ReadFromFile(fileName)
}

func (p *PgsStorage) SaveToFile(fileName string) {
	p.memory.SaveToFile(fileName)
}