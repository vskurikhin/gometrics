/*
 * This file was last modified at 2024-06-16 14:10 by Victor N. Skurikhin.
 * config_test.go
 * $Id$
 */

package env

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	testDataBaseDSN   = "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	testKey           string
	testServerAddress string
	testTempFileName  string
)

func TestConfig(t *testing.T) {

	getEnvironments()
	initServerFlags()
	reportInterval := 1 * time.Nanosecond
	flag.reportInterval = &reportInterval
	pollInterval := 2 * time.Nanosecond
	flag.pollInterval = &pollInterval
	getTestConfig()

	var tests = []struct {
		name string
		fCfg func() Config
		want string
	}{
		{
			name: "Test config #1",
			fCfg: GetAgentConfig,
			want: `
	dataBaseDSN     : 
	fileStoragePath : 
	key             : 
	pollInterval    : 2ns
	reportInterval  : 1ns
	restore         : false
	serverAddress   : localhost:8080
	storeInterval   : 0s
	urlHost         : http://localhost:8080
`,
		},
		{
			name: "Test config #2",
			fCfg: GetServerConfig,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : /tmp/metrics-db.json
	key             : 
	pollInterval    : 0s
	reportInterval  : 0s
	restore         : true
	serverAddress   : localhost:8080
	storeInterval   : 5m0s
	urlHost         : http://localhost:8080
`,
		},
		{
			name: "Test config #3",
			fCfg: getTestConfig,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + testTempFileName + `
	key             : ` + testKey + `
	pollInterval    : 1m0s
	reportInterval  : 1h0m0s
	restore         : true
	serverAddress   : ` + testServerAddress + `
	storeInterval   : 24h0m0s
	urlHost         : http://` + testServerAddress + `
`,
		},
		{
			name: "Test config #4",
			fCfg: func() Config { return getTestEnvironments(t) },
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + testTempFileName + `
	key             : ` + testKey + `
	pollInterval    : 0s
	reportInterval  : 0s
	restore         : true
	serverAddress   : ` + testServerAddress + `
	storeInterval   : 1s
	urlHost         : http://` + testServerAddress + `
`,
		},
		{
			name: "Test config #5",
			fCfg: getTestFlags,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + testTempFileName + `
	key             : ` + testKey + `
	pollInterval    : 1m0s
	reportInterval  : 1h0m0s
	restore         : true
	serverAddress   : ` + testServerAddress + `
	storeInterval   : 24h0m0s
	urlHost         : http://` + testServerAddress + `
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			onceCfg = new(sync.Once)
			c := test.fCfg()
			got := c.String()
			assert.Equal(t, test.want, got)
		})
	}
	assert.True(t, cfg.IsDBSetup())
}

func getTestConfig() Config {
	return GetTestConfig(
		WithDataBaseDSN(&testDataBaseDSN),
		WithFileStoragePath(testTempFileName),
		WithKey(&testKey),
		WithPollInterval(time.Minute),
		WithReportInterval(time.Hour),
		WithRestore(true),
		WithServerAddress(testServerAddress),
		WithStoreInterval(24*time.Hour),
	)
}

func getTestEnvironments(t *testing.T) Config {

	onceEnv = new(sync.Once)

	t.Setenv("ADDRESS", testServerAddress)
	t.Setenv("DATABASE_DSN", testDataBaseDSN)
	t.Setenv("FILE_STORAGE_PATH", testTempFileName)
	t.Setenv("KEY", testKey)
	t.Setenv("RESTORE", "true")
	t.Setenv("STORE_INTERVAL", "1")

	return GetServerConfig()
}

func getTestFlags() Config {
	onceFlags = new(sync.Once)
	initAgentFlags()
	return GetTestConfig(
		WithDataBaseDSN(&testDataBaseDSN),
		WithFileStoragePath(testTempFileName),
		WithKey(&testKey),
		WithPollInterval(time.Minute),
		WithReportInterval(time.Hour),
		WithRestore(true),
		WithServerAddress(testServerAddress),
		WithStoreInterval(24*time.Hour),
	)
}

func init() {
	port := 65500 + rand.Intn(34)
	testKey = fmt.Sprintf("%018d", rand.Uint32())
	testServerAddress = fmt.Sprintf("localhost:%d", port)
	testTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}
