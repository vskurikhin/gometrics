/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * reports_test.go
 * $Id$
 */

package agent

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/vskurikhin/gometrics/internal/env"
	t "github.com/vskurikhin/gometrics/internal/types"
)

var (
	enabled           = []t.Name{t.TotalAlloc, t.PollCount, t.RandomValue}
	testDataBaseDSN   = ""
	testKey           string
	testServerAddress string
	testTempFileName  string
)

func TestReports(t *testing.T) {
	s := "1"
	Storage()
	store.Put("PollCount", &s)
	store.Put("RandomValue", &s)
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		_, _ = res.Write([]byte(""))
	}))

	a := strings.Split(testServer.URL, "://")
	if len(a) < 2 {
		t.Fatalf("len(%s) < 2", a)
	}
	t.Setenv("REPORT_INTERVAL", "1")
	t.Setenv("ADDRESS", a[1])
	cfg := getTestConfig()
	client := http.Client{}
	reports(cfg, enabled, &client)
	testServer.Close()
	reports(cfg, enabled, &client)
}

func TestIsUpperBound(t *testing.T) {
	var d, i, s = 5, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = (1 << i)
	}
	assert.True(t, s < d)
	d, i, s = 10, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
	}
	assert.True(t, (1<<i) < d)
	d, i, s = 25, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
	}
	assert.True(t, (1<<i) < d)
	d, i, s = 50, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
	}
	assert.True(t, (1<<i) < d)
}

func getTestConfig() env.Config {
	return env.GetTestConfig(
		env.GetProperty,
		env.WithDataBaseDSN(&testDataBaseDSN),
		env.WithFileStoragePath(testTempFileName),
		env.WithKey(&testKey),
		env.WithPollInterval(30*time.Minute),
		env.WithReportInterval(1),
		env.WithRestore(true),
		env.WithServerAddress(testServerAddress),
		env.WithStoreInterval(24*time.Hour),
		env.WithTrustedSubnet("127.0.0.0/8"),
	)
}

func init() {
	port := 65500 + rand.Intn(34)
	testKey = fmt.Sprintf("%018d", rand.Uint32())
	testServerAddress = fmt.Sprintf("localhost:%d", port)
	testTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}
