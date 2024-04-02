/*
 * This file was last modified at 2024-03-19 12:12 by Victor N. Skurikhin.
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
	router.Get(names.Ping, handlers.PingHandler)
	router.Post(names.UpdateChi, handlers.UpdateHandler)
	router.Post(names.UpdateURL, handlers.UpdateJSONHandler)
	router.Post(names.UpdatesURL, handlers.UpdatesJSONHandler)
	router.Get(names.ValueChi, handlers.ValueHandler)
	router.Post(names.ValueURL, handlers.ValueJSONHandler)

	return router
}
