/*
 * This file was last modified at 2024-03-18 23:39 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"net/http"
)

func main() {

	env.InitServer()
	if env.Server.IsDBSetup() {
		server.DBConnect()
		server.CreateSchema()
	}
	server.Storage()
	server.Read()

	r := chi.NewRouter()
	r.Use(compress.Compress)
	r.Use(logger.Logging)
	r.Use(middleware.Recoverer)
	r.Get("/", handlers.RootHandler)
	r.Get(names.Ping, handlers.PingHandler)
	r.Post(names.UpdateChi, handlers.UpdateHandler)
	r.Post(names.UpdateURL, handlers.UpdateJSONHandler)
	r.Post(names.UpdatesURL, handlers.UpdatesJSONHandler)
	r.Get(names.ValueChi, handlers.ValueHandler)
	r.Post(names.ValueURL, handlers.ValueJSONHandler)

	go server.Save()
	go server.DBPing()
	err := http.ListenAndServe(env.Server.ServerAddress(), r)
	if err != nil {
		panic(err)
	}
}
