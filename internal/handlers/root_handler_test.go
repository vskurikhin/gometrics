/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * root_handler_test.go
 * $Id$
 */

package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/storage"
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
	env.GetServerConfig()
	mem := new(storage.MemStorage)
	mem.Metrics = make(map[string]*string)
	store = mem
	return env.GetTestConfig(
		func() env.Property {
			return env.GetTestProperty(env.WithStorage(store))
		},
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
