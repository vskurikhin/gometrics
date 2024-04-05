/*
 * This file was last modified at 2024-04-03 08:47 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"net/http"
)

func main() {

	env.InitServer()
	server.DBInit()
	server.Storage()
	server.Read()

	router := initRouter()

	go server.Save()
	err := http.ListenAndServe(env.Server.ServerAddress(), router)
	if err != nil {
		panic(err)
	}
}

func initRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Use(compress.Compress)
	router.Use(logger.Logging)
	router.Use(middleware.Recoverer)
	router.Get("/", handlers.RootHandler)
	router.Get(env.Ping, handlers.PingHandler)
	router.Post(env.UpdateChi, handlers.UpdateHandler)
	router.Post(env.UpdateURL, handlers.UpdateJSONHandler)
	router.Post(env.UpdatesURL, handlers.UpdatesJSONHandler)
	router.Get(env.ValueChi, handlers.ValueHandler)
	router.Post(env.ValueURL, handlers.ValueJSONHandler)

	return router
}
