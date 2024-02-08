/*
 * This file was last modified at 2024-02-04 12:18 by Victor N. Skurikhin.
 * main.go
 * $Id$
 */

package main

import (
	"fmt"
	"github.com/vskurikhin/gometrics/api/names"
	"github.com/vskurikhin/gometrics/internal/handlers"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle(names.UpdateURL, http.HandlerFunc(handlers.UpdateHandler))

	err := http.ListenAndServe(fmt.Sprintf(":%d", names.Port), mux)
	if err != nil {
		panic(err)
	}
}
