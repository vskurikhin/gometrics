/*
 * This file was last modified at 2024-06-25 00:37 by Victor N. Skurikhin.
 * main_test.go
 * $Id$
 */

package main

import (
	"context"
	"fmt"
	"math/rand"
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

func TestRun(t *testing.T) {

	t.Setenv("ADDRESS", testServerAddress)
	t.Setenv("DATABASE_DSN", testDataBaseDSN)
	t.Setenv("REPORT_INTERVAL", "1")
	t.Setenv("POLL_INTERVAL", "1")
	t.Setenv("KEY", testKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	run(ctx)
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func init() {
	port := 65500 + rand.Intn(34)
	testKey = fmt.Sprintf("%018d", rand.Uint32())
	testServerAddress = fmt.Sprintf("localhost:%d", port)
	testTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}
