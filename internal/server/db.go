/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * db.go
 * $Id$
 */

package server

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/vskurikhin/gometrics/internal/env"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vskurikhin/gometrics/internal/logger"
)

type DBHealth interface {
	GetStatus() bool
	checkStatus() error
}

type PgxPool interface {
	getPool() *pgxpool.Pool
}

type PgxPoolHealth struct {
	sync.RWMutex
	pool   *pgxpool.Pool
	status bool
}

var dbHealth = new(PgxPoolHealth)

func DBInit(cfg env.Config) {
	dbHealth.pool = cfg.Property().DBPool()
	if cfg.IsDBSetup() {
		CreateSchema()
		go dbPing()
	}
}

func pgxPoolInstance() PgxPool {
	return dbHealth
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

func (p *PgxPoolHealth) getPool() *pgxpool.Pool {
	return dbHealth.pool
}

func dbPing() {
	for {
		time.Sleep(2 * time.Second)
		err := dbHealth.checkStatus()
		if err != nil {
			logger.Log.Debug("db health checkStatus ", zap.String("error", fmt.Sprintf("%v", err)))
		}
	}
}
