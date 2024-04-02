/*
 * This file was last modified at 2024-03-19 09:58 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/parser"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/types"
	"net/http"
	"os"
)

func UpdateHandler(response http.ResponseWriter, request *http.Request) {

	store = server.Storage()

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
		store.PutCounter(name, parsed.CalcValue(store.GetCounter(name)))
	case types.GAUGE:
		store.PutGauge(name, parsed.CalcValue(store.GetGauge(name)))
	}
	response.WriteHeader(http.StatusOK)
}
