/*
 * This file was last modified at 2024-02-10 15:27 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/parser"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"net/http"
	"os"
)

var storage = memory.Instance()

func UpdateHandler(response http.ResponseWriter, request *http.Request) {

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
	storage.Put(name, parsed.CalcValue(storage.Get(name)))
}
