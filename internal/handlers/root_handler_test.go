/*
 * This file was last modified at 2024-06-11 12:34 by Victor N. Skurikhin.
 * root_handler_test.go
 * $Id$
 */

package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vskurikhin/gometrics/internal/env"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var (
	testDataBaseDSN   = ""
	testKey           string
	testServerAddress string
	testTempFileName  string
)

func TestRootHandler(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/type/PollCount", nil)

	ctx := chi.NewRouteContext()

	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	RootHandler(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func getTestConfig() env.Config {
	return env.GetTestConfig(
		env.WithDataBaseDSN(&testDataBaseDSN),
		env.WithFileStoragePath(testTempFileName),
		env.WithKey(&testKey),
		env.WithPollInterval(30*time.Minute),
		env.WithReportInterval(time.Hour),
		env.WithRestore(true),
		env.WithServerAddress(testServerAddress),
		env.WithStoreInterval(24*time.Hour),
	)
}

func init() {
	port := 65500 + rand.Intn(34)
	testKey = fmt.Sprintf("%018d", rand.Uint32())
	testServerAddress = fmt.Sprintf("localhost:%d", port)
	testTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}
