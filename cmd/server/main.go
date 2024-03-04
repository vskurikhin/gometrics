/*
 * This file was last modified at 2024-03-02 20:31 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/compress"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/server"
	"go.uber.org/zap"
	"net/http"
)

func main() {

	env.InitServer()
	logger.Log.Debug("Server ", zap.String("env", fmt.Sprintf("%+v", env.Server)))
	server.Read()

	r := chi.NewRouter()
	r.Use(compress.Compress)
	r.Use(logger.Logging)
	r.Use(middleware.Recoverer)
	r.Get("/", handlers.RootHandler)
	r.Post(names.UpdateChi, handlers.UpdateHandler)
	r.Post(names.UpdateURL, handlers.UpdateJSONHandler)
	r.Get(names.ValueChi, handlers.ValueHandler)
	r.Post(names.ValueURL, handlers.ValueJSONHandler)

	go server.Save()
	err := http.ListenAndServe(env.Server.ServerAddress(), r)
	if err != nil {
		panic(err)
	}
}
