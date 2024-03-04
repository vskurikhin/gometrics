/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/parser"
	"net/http"
	"os"
)

func UpdateHandler(response http.ResponseWriter, request *http.Request) {

	logger.Log.Debug("got incoming HTTP request with Text")

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
	store.Put(name, parsed.CalcValue(store.Get(name)))
	response.WriteHeader(http.StatusOK)
}
