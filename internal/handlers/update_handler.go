/*
 * This file was last modified at 2024-03-18 19:14 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/parser"
	"github.com/vskurikhin/gometrics/internal/storage/postgres"
	"github.com/vskurikhin/gometrics/internal/types"
	"net/http"
	"os"
)

func UpdateHandler(response http.ResponseWriter, request *http.Request) {

	store = postgres.Instance()

	defer func() {
		if p := recover(); p != nil {
			//goland:noinspection GoUnhandledErrorResult
			fmt.Fprintf(os.Stderr, "update error: %v", p)
			response.WriteHeader(http.StatusNotFound)
		}
	}()

	parsed, err := parser.Parse(request)
	if err != nil || parsed.Value() == nil {
		response.WriteHeader(parsed.Status())
		return
	}

	name := parsed.String()
	switch parsed.Type() {
	case types.COUNTER:
		store.PutCounter(name, parsed.CalcValue(store.Get(name)))
	case types.GAUGE:
		store.PutGauge(name, parsed.CalcValue(store.Get(name)))
	}
	response.WriteHeader(http.StatusOK)
}
