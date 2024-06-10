/*
 * This file was last modified at 2024-06-10 21:53 by Victor N. Skurikhin.
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
	TestDataBaseDSN   = "postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable"
	TestKey           string
	TestServerAddress string
	TestTempFileName  string
)

func init() {
	port := 65500 + rand.Intn(99)
	TestKey = fmt.Sprintf("%018d", rand.Uint32())
	TestServerAddress = fmt.Sprintf("localhost:%d", port)
	TestTempFileName = fmt.Sprintf("%s/test_%018d.txt", os.TempDir(), rand.Uint32())
}

func GetTestConfig(opts ...func(*config)) Config {
	onceCfg.Do(func() {
		cfg = getConfig(opts...)
	})
	return cfg
}

func getTestConfig() Config {
	return GetTestConfig(
		WithDataBaseDSN(&TestDataBaseDSN),
		WithFileStoragePath(TestTempFileName),
		WithKey(&TestKey),
		WithPollInterval(time.Minute),
		WithReportInterval(time.Hour),
		WithRestore(true),
		WithServerAddress(TestServerAddress),
		WithStoreInterval(24*time.Hour),
	)
}

func getTestEnvironments(t *testing.T) Config {

	onceEnv = new(sync.Once)

	t.Setenv("ADDRESS", TestServerAddress)
	t.Setenv("STORE_INTERVAL", "1")
	t.Setenv("FILE_STORAGE_PATH", TestTempFileName)
	t.Setenv("RESTORE", "true")
	t.Setenv("DATABASE_DSN", TestDataBaseDSN)
	t.Setenv("KEY", TestKey)

	return GetServerConfig()
}

func getTestFlags() Config {
	onceFlags = new(sync.Once)
	initAgentFlags()
	return GetTestConfig(
		WithDataBaseDSN(&TestDataBaseDSN),
		WithFileStoragePath(TestTempFileName),
		WithKey(&TestKey),
		WithPollInterval(time.Minute),
		WithReportInterval(time.Hour),
		WithRestore(true),
		WithServerAddress(TestServerAddress),
		WithStoreInterval(24*time.Hour),
	)
}

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
	storeInterval   : 300ns
	urlHost         : http://localhost:8080
`,
		},
		{
			name: "Test config #3",
			fCfg: getTestConfig,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + TestTempFileName + `
	key             : ` + TestKey + `
	pollInterval    : 1m0s
	reportInterval  : 1h0m0s
	restore         : true
	serverAddress   : ` + TestServerAddress + `
	storeInterval   : 24h0m0s
	urlHost         : http://` + TestServerAddress + `
`,
		},
		{
			name: "Test config #4",
			fCfg: func() Config { return getTestEnvironments(t) },
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + TestTempFileName + `
	key             : ` + TestKey + `
	pollInterval    : 0s
	reportInterval  : 0s
	restore         : true
	serverAddress   : ` + TestServerAddress + `
	storeInterval   : 1ns
	urlHost         : http://` + TestServerAddress + `
`,
		},
		{
			name: "Test config #5",
			fCfg: getTestFlags,
			want: `
	dataBaseDSN     : postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable
	fileStoragePath : ` + TestTempFileName + `
	key             : ` + TestKey + `
	pollInterval    : 1m0s
	reportInterval  : 1h0m0s
	restore         : true
	serverAddress   : ` + TestServerAddress + `
	storeInterval   : 24h0m0s
	urlHost         : http://` + TestServerAddress + `
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
