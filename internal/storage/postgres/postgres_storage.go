/*
 * This file was last modified at 2024-03-18 19:23 by Victor N. Skurikhin.
 * postgres_storage.go
 * $Id$
 */

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

//goland:noinspection GoNameStartsWithPackageName
type PostgresStorage struct {
	memory storage.Storage
	pool   *pgxpool.Pool
}

var initMutex sync.Mutex
var postgres *PostgresStorage

func Instance() storage.Storage {

	initMutex.Lock()
	defer initMutex.Unlock()

	if postgres == nil {
		postgres = new(PostgresStorage)
		postgres.pool = server.
			PgxPoolInstance().
			GetPool()
		postgres.memory = memory.Instance()
	}
	return postgres
}

func (p *PostgresStorage) Get(name string) *string {
	return p.memory.Get(name)
}

func (p *PostgresStorage) GetCounter(name string) *string {

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func GetCounter", zap.String("error", fmt.Sprintf("%v", p)))
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
		return p.memory.Get(name)
	}
	row := conn.QueryRow(
		ctx,
		"SELECT counter FROM metric WHERE name = $1 AND type = 'counter'",
		name,
	)
	var counter sql.NullInt64

	err = row.Scan(&counter)

	if err != nil {
		panic(err)
	}

	if counter.Valid {
		result := fmt.Sprintf("%d", counter.Int64)
		return &result
	}
	return nil
}

func (p *PostgresStorage) GetGauge(name string) *string {

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func GetGauge", zap.String("error", fmt.Sprintf("%v", p)))
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
		return p.memory.Get(name)
	}
	row := conn.QueryRow(
		ctx,
		"SELECT gauge FROM metric WHERE name = $1 AND type = 'gauge'",
		name,
	)
	var gauge sql.NullFloat64

	err = row.Scan(&gauge)

	if err != nil {
		panic(err)
	}

	if gauge.Valid {
		result := fmt.Sprintf("%.3f", gauge.Float64)
		if result[len(result)-1] == '0' {
			result = result[:len(result)-1]
		}
		return &result
	}
	return nil
}

func (p *PostgresStorage) Put(name string, value *string) {
	p.memory.Put(name, value)
}

func (p *PostgresStorage) PutCounter(name string, value *string) {

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func PutCounter", zap.String("error", fmt.Sprintf("%v", p)))
		}
	}()

	p.memory.Put(name, value)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, _ := p.pool.Acquire(ctx)
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	_, _ = conn.Exec(ctx, "INSERT INTO name (name) VALUES ($1) ON CONFLICT DO NOTHING", name)

	counter, _ := strconv.Atoi(*value)
	_, _ = conn.Exec(ctx,
		`INSERT INTO metric (name, type, counter)
				VALUES ($1, 'counter', $2)
				ON CONFLICT(name) 
				DO UPDATE SET
				counter = $2`,
		name, counter,
	)
}

func (p *PostgresStorage) PutGauge(name string, value *string) {

	defer func() {
		if p := recover(); p != nil {
			logger.Log.Debug("func PutGauge", zap.String("error", fmt.Sprintf("%v", p)))
		}
	}()

	p.memory.Put(name, value)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, _ := p.pool.Acquire(ctx)
	defer func() {
		if conn != nil {
			conn.Release()
		}
	}()

	_, _ = conn.Exec(ctx, "INSERT INTO name (name) VALUES ($1) ON CONFLICT DO NOTHING", name)

	gauge, _ := strconv.ParseFloat(*value, 64)
	_, _ = conn.Exec(ctx,
		`INSERT INTO metric (name, type, gauge)
				VALUES ($1, 'gauge', $2)
				ON CONFLICT(name) 
				DO UPDATE SET
				gauge = $2`,
		name, gauge,
	)
}

func (p *PostgresStorage) ReadFromFile(fileName string) {
	p.memory.ReadFromFile(fileName)
}

func (p *PostgresStorage) SaveToFile(fileName string) {
	p.memory.SaveToFile(fileName)
}
