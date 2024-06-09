/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * main_test.go.go
 * $Id$
 */

package main

import (
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/server"
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

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	router.Mount("/debug", middleware.Profiler())

	router.Group(func(r chi.Router) {
		r.Post(env.UpdateChi, handlers.UpdateHandler)
		r.Get(env.ValueChi, handlers.ValueHandler)
		r.Post(env.UpdatesURL, handlers.UpdatesJSONHandler)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.Compress(9))
		r.Get("/", handlers.RootHandler)
		r.Post(env.UpdateURL, handlers.UpdateJSONHandler)
		r.Post(env.ValueURL, handlers.ValueJSONHandler)
	})

	return router
}
