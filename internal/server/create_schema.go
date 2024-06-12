/*
 * This file was last modified at 2024-03-18 22:41 by Victor N. Skurikhin.
 * create_schema.go
 * $Id$
 */

// Package server реализация серверных частей
package server

import (
	"context"
	"time"
)

var pgxPool = PgxPoolInstance()

func CreateSchema() {
	pool := pgxPool.GetPool()
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
	_, err = conn.Exec(ctx, "CREATE TYPE public.TYPE AS ENUM ('gauge', 'counter')")
	if err != nil {
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
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(ctx,
		`CREATE TABLE IF NOT EXISTS name (
					name TEXT NOT NULL UNIQUE
					)`,
	)
	if err != nil {
		panic(err)
	}
}
