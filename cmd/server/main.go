/*
 * This file was last modified at 2024-03-01 21:41 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/logger"
	"net/http"
)

func main() {

	env.InitServer()

	r := chi.NewRouter()
	r.Use(compress.Compress)
	r.Use(logger.Logging)
	r.Use(middleware.Recoverer)
	r.Get("/", handlers.RootHandler)
	r.Post(names.UpdateChi, handlers.UpdateHandler)
	r.Post(names.UpdateURL, handlers.UpdateJSONHandler)
	r.Get(names.ValueChi, handlers.ValueHandler)
	r.Post(names.ValueURL, handlers.ValueJSONHandler)

	err := http.ListenAndServe(env.Server.ServerAddress(), r)
	if err != nil {
		panic(err)
	}
}
