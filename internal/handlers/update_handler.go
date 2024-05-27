/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vskurikhin/gometrics/internal/parser"
	"github.com/vskurikhin/gometrics/internal/server"
	"github.com/vskurikhin/gometrics/internal/types"
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
