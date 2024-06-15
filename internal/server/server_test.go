/*
 * This file was last modified at 2024-06-16 14:35 by Victor N. Skurikhin.
 * server_test.go
 * $Id$
 */

package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gometrics/internal/env"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	dbInitDataBaseDSN = "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	testDataBaseDSN   = ""
	testKey           string
	testServerAddress string
	testTempFileName  string
)

func TestDBInit(t *testing.T) {

	cfg := getTestDBInit()
	DBInit(cfg)
	time.Sleep(2 * time.Second)
	assert.NotNil(t, pgxPool.getPool())
}

//TODO
//func TestSave(t *testing.T) {
//
//	cfg := getTestConfig()
//	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
//	defer cancel()
//	Storage(cfg)
//	Save(ctx, cfg)
//	Read(cfg)
//}

func TestServer(t *testing.T) {

	getTestConfig()

	var tests = []struct {
		name string
	}{
		{
			name: "Test config #1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//got := c.String()
			//assert.Equal(t, test.want, got)
		})
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

func getTestDBInit() env.Config {
	return env.GetTestConfig(
		env.WithDataBaseDSN(&dbInitDataBaseDSN),
		env.WithFileStoragePath(testTempFileName),
		env.WithKey(&testKey),
		env.WithPollInterval(30*time.Minute),
		env.WithReportInterval(time.Hour),
		env.WithRestore(false),
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
