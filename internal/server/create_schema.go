/*
 * This file was last modified at 2024-06-11 09:46 by Victor N. Skurikhin.
 * create_schema.go
 * $Id$
 */

// Package server реализация серверных частей
package server

import (
	"context"
	"github.com/vskurikhin/gometrics/internal/util"
	"time"
)

var pgxPool = pgxPoolInstance()

func CreateSchema() {
	pool := pgxPool.getPool()
	if pool == nil {
		panic("poll is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer func() {
		cancel()
		ctx.Done()
	}()

	conn, err := pool.Acquire(ctx)
	defer conn.Release()

	if conn == nil || err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS metric (
					id BIGSERIAL,
					name TEXT NOT NULL UNIQUE,
					type public.TYPE NOT NULL,
					gauge DOUBLE PRECISION,
					counter BIGINT
					)`,
	)
	util.IfErrorThenPanic(err)
	_, err = conn.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS name (
					name TEXT NOT NULL UNIQUE
					)`,
	)
	util.IfErrorThenPanic(err)
}
