/*
 * This file was last modified at 2024-03-18 19:27 by Victor N. Skurikhin.
 * db.go
 * $Id$
 */

package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"go.uber.org/zap"
	"sync"
	"time"
)

type DBHealth interface {
	GetStatus() bool
	checkStatus() error
}

type PgxPool interface {
	GetPool() *pgxpool.Pool
}

type PgxPoolHealth struct {
	sync.RWMutex
	pool   *pgxpool.Pool
	status bool
}

const tryLockCheckStatus = 1000

var dbHealth = new(PgxPoolHealth)

func DBHealthInstance() DBHealth {
	return dbHealth
}

func PgxPoolInstance() PgxPool {
	return dbHealth
}

func DBConnect() *pgxpool.Pool {

	config, err := pgxpool.ParseConfig(env.Server.DataBaseDSN())
	if err != nil {
		panic(err)
	}
	logger.Log.Debug("DBConnect config parsed")

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		logger.Log.Debug("Acquire connect ping...")
		if err = conn.Ping(ctx); err != nil {
			panic(err)
		}
		logger.Log.Debug("Acquire connect Ok")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), config)
	if err != nil {
		panic(err)
	}
	logger.Log.Debug("NewWithConfig pool created")
	_, err = pool.Acquire(context.TODO())
	if err != nil {
		panic(err)
	}
	logger.Log.Debug("Acquire pool Ok")
	dbHealth.pool = pool

	return pool
}

func (p *PgxPoolHealth) GetStatus() bool {
	p.RLock()
	defer p.RUnlock()
	return dbHealth.status
}

func (p *PgxPoolHealth) checkStatus() error {

	p.Lock()
	defer p.Unlock()

	if p.pool == nil {
		p.status = false
		return errors.New("poll is nil")
	}

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
		p.status = false
		return err
	}
	p.status = true

	return nil
}

func (p *PgxPoolHealth) GetPool() *pgxpool.Pool {
	return dbHealth.pool
}

func DBPing() {
	for {
		dbPing()
	}
}

func dbPing() {
	time.Sleep(2 * time.Second)
	err := dbHealth.checkStatus()
	if err != nil {
		logger.Log.Debug("db health checkStatus ", zap.String("error", fmt.Sprintf("%v", err)))
	}
}
