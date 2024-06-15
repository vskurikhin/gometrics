/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/render"
	"github.com/vskurikhin/gometrics/internal/util"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/server"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	ctx := context.Background()
	fmt.Printf(
		"Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildVersion, buildDate, buildCommit,
	)
	cfg := env.GetServerConfig()
	fmt.Print(cfg)

	server.DBInit(cfg)
	server.Storage(cfg)
	server.Read(cfg)

	router := initRouter()

	go func() {
		err := http.ListenAndServe(cfg.ServerAddress(), router)
		util.IfErrorThenPanic(err)
	}()
	server.Save(ctx, cfg)
}

func initRouter() *chi.Mux {

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	router.Mount("/debug", middleware.Profiler())

	router.Group(func(r chi.Router) {
		r.Post(env.UpdateChi, handlers.UpdateHandler)
		r.Get(env.ValueChi, handlers.ValueHandler)
	})

	router.Group(func(r chi.Router) {
		r.Use(middleware.Compress(9))
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Get("/", handlers.RootHandler)
		r.Post(env.UpdateURL, handlers.UpdateJSONHandler)
		r.Post(env.ValueURL, handlers.ValueJSONHandler)
	})

	router.Group(func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))
		r.Post(env.UpdatesURL, handlers.UpdatesJSONHandler)
	})

	return router
}
