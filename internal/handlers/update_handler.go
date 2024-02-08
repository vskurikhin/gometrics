/*
 * This file was last modified at 2024-02-04 13:29 by Victor N. Skurikhin.
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

func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		defer func() {
			if p := recover(); p != nil {
				//goland:noinspection GoUnhandledErrorResult
				fmt.Fprintf(os.Stderr, "update error: %v", p)
				w.WriteHeader(http.StatusNotFound)
			}
		}()

		parsed, err := parser.Parse(r)
		if err != nil || parsed.Value() == nil {
			w.WriteHeader(parsed.Status())
			return
		}
		name := parsed.String()
		storage := memory.Instance()
		storage.Put(name, *parsed.Value())

		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
