/*
 * This file was last modified at 2024-06-11 12:34 by Victor N. Skurikhin.
 * reports_test.go
 * $Id$
 */

package agent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gometrics/internal/env"
	t "github.com/vskurikhin/gometrics/internal/types"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
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
		res.Write([]byte(""))
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
	var d, i, s int = 5, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
		fmt.Fprintf(os.Stderr, "d: %d, i: %d, s: %d, isUpperBound: %v\n", time.Duration(d)*time.Second, i, s, isUpperBound(i, time.Duration(d)))
	}
	assert.True(t, s < d)
	fmt.Fprintln(os.Stderr)
	d, i, s = 10, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
		fmt.Fprintf(os.Stderr, "d: %d, i: %d, s: %d, isUpperBound: %v\n", time.Duration(d)*time.Second, i, s, isUpperBound(i, time.Duration(d)))
	}
	assert.True(t, (1<<i) < d)
	fmt.Fprintln(os.Stderr)
	d, i, s = 25, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
		fmt.Fprintf(os.Stderr, "d: %d, i: %d, s: %d, isUpperBound: %v\n", time.Duration(d)*time.Second, i, s, isUpperBound(i, time.Duration(d)))
	}
	assert.True(t, (1<<i) < d)
	fmt.Fprintln(os.Stderr)
	d, i, s = 50, 0, 0
	for ; isUpperBound(i, time.Duration(d)); i++ {
		s = s + (1 << i)
		fmt.Fprintf(os.Stderr, "d: %d, i: %d, s: %d, isUpperBound: %v\n", time.Duration(d)*time.Second, i, s, isUpperBound(i, time.Duration(d)))
	}
	assert.True(t, (1<<i) < d)
}

func getTestConfig() env.Config {
	return env.GetTestConfig(
		env.WithDataBaseDSN(&testDataBaseDSN),
		env.WithFileStoragePath(testTempFileName),
		env.WithKey(&testKey),
		env.WithPollInterval(30*time.Minute),
		env.WithReportInterval(1),
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
