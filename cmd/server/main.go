/*
 * This file was last modified at 2024-02-24 17:37 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"github.com/vskurikhin/gometrics/internal/logger"
	"net/http"
)

func main() {

	env.InitServer()

	r := chi.NewRouter()
	r.Use(logger.WithLogging)
	r.Post(names.UpdateChi, handlers.UpdateHandler)
	r.Get(names.ValueChi, handlers.ValueHandler)

	err := http.ListenAndServe(env.Server.ServerAddress(), r)
	if err != nil {
		panic(err)
	}
}
