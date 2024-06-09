/*
 * This file was last modified at 2024-06-11 12:34 by Victor N. Skurikhin.
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
	t.Setenv("FILE_STORAGE_PATH", testTempFileName)
	t.Setenv("KEY", testKey)
	t.Setenv("RESTORE", "true")
	t.Setenv("STORE_INTERVAL", "1")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()
	run(ctx)
}

func TestM(t *testing.T) {
}

func init() {
	port := 65500 + rand.Intn(34)
	testKey = fmt.Sprintf("%018d", rand.Uint32())
	testServerAddress = fmt.Sprintf("localhost:%d", port)
	testTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}
