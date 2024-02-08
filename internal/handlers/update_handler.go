/*
 * This file was last modified at 2024-02-08 21:41 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/parser"
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"github.com/vskurikhin/gometrics/internal/util"
	"net/http"
	"os"
)

var storage = memory.Instance()

func UpdateHandler(w http.ResponseWriter, r *http.Request) {

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

	var name string
	if parsed.Name() > 0 {
		name = parsed.String()

	} else {
		path := util.SplitPath(r)
		name = path[2]
	}
	storage.Put(name, parsed.CalcValue(storage.Get(name)))
}
