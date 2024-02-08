/*
 * This file was last modified at 2024-02-04 12:18 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Post(names.UpdateChi, handlers.UpdateHandler)
	r.Get(names.ValueChi, handlers.ValueHandler)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
