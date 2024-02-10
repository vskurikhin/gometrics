/*
 * This file was last modified at 2024-02-10 15:07 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/spf13/pflag"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/cmd/cflag"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"net/http"
)

func main() {

	pflag.Parse()

	r := chi.NewRouter()
	r.Post(names.UpdateChi, handlers.UpdateHandler)
	r.Get(names.ValueChi, handlers.ValueHandler)

	err := http.ListenAndServe(cflag.ServerFlags.ServerAddress(), r)
	if err != nil {
		panic(err)
	}
}
