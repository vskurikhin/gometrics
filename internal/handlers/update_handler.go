/*
 * This file was last modified at 2024-02-04 13:29 by Victor N. Skurikhin.
 * update_handler.go
 * $Id$
 */

package handlers

import (
	"github.com/vskurikhin/gometrics/internal/storage/memory"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		parsed, err := parse(r)
		if err != nil {
			w.WriteHeader(parsed.status)
		}
		storage := memory.Instance()
		_ = storage.Put(parsed.n, parsed.value)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
