/*
 * This file was last modified at 2024-06-10 18:44 by Victor N. Skurikhin.
 * config.go
 * $Id$
 */

// Package env работа с настройками и окружением
package env

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/util"
	"strconv"
	"sync"
	"time"
)

type Config interface {
	fmt.Stringer
	DataBaseDSN() string
	FileStoragePath() string
	IsDBSetup() bool
	Key() *string
	PollInterval() time.Duration
	ReportInterval() time.Duration
	Restore() bool
	ServerAddress() string
	StoreInterval() time.Duration
	URLHost() *string
}

type config struct {
	dataBaseDSN     *string
	fileStoragePath string
	key             *string
	pollInterval    time.Duration
	reportInterval  time.Duration
	restore         bool
	serverAddress   string
	storeInterval   time.Duration
	urlHost         *string
}

var onceCfg = new(sync.Once)
var cfg *config

func GetAgentConfig() Config {
	onceCfg.Do(func() {
		cfg = new(config)
		getEnvironments()
		initAgentFlags()
		initAgentConfig()
	})
	return cfg
}

func GetServerConfig() Config {
	onceCfg.Do(func() {
		cfg = new(config)
		getEnvironments()
		initServerFlags()
		initServerConfig()
	})
	return cfg
}

// WithDataBaseDSN — пример функции, которая присваивает поле Mode.
func WithDataBaseDSN(dataBaseDSN *string) func(*config) {
	return func(c *config) {
		c.dataBaseDSN = dataBaseDSN
	}
}

func (c *config) DataBaseDSN() string {
	if c.dataBaseDSN != nil {
		return *c.dataBaseDSN
	}
	return ""
}

// WithFileStoragePath — пример функции, которая присваивает поле Mode.
func WithFileStoragePath(fileStoragePath string) func(*config) {
	return func(c *config) {
		c.fileStoragePath = fileStoragePath
	}
}

func (c *config) FileStoragePath() string {
	return c.fileStoragePath
}

func (c *config) IsDBSetup() bool {
	return c.dataBaseDSN != nil && *c.dataBaseDSN != ""
}

// WithKey — пример функции, которая присваивает поле Mode.
func WithKey(key *string) func(*config) {
	return func(c *config) {
		c.key = key
	}
}

func (c *config) Key() *string {
	return c.key
}

// WithPollInterval — пример функции, которая присваивает поле Mode.
func WithPollInterval(pollInterval time.Duration) func(*config) {
	return func(c *config) {
		c.pollInterval = pollInterval
	}
}

func (c *config) PollInterval() time.Duration {
	return c.pollInterval
}

// WithReportInterval — пример функции, которая присваивает поле Mode.
func WithReportInterval(reportInterval time.Duration) func(*config) {
	return func(c *config) {
		c.reportInterval = reportInterval
	}
}

func (c *config) ReportInterval() time.Duration {
	return c.reportInterval
}

// WithRestore — пример функции, которая присваивает поле Mode.
func WithRestore(restore bool) func(*config) {
	return func(c *config) {
		c.restore = restore
	}
}

func (c *config) Restore() bool {
	return c.restore
}

// WithServerAddress — пример функции, которая присваивает поле Mode.
func WithServerAddress(serverAddress string) func(*config) {
	return func(c *config) {
		c.serverAddress = serverAddress
	}
}

func (c *config) ServerAddress() string {
	return c.serverAddress
}

// WithStoreInterval — пример функции, которая присваивает поле Mode.
func WithStoreInterval(storeInterval time.Duration) func(*config) {
	return func(c *config) {
		c.storeInterval = storeInterval
	}
}

func (c *config) StoreInterval() time.Duration {
	return c.storeInterval
}

func (c *config) String() string {
	format := `
	dataBaseDSN     : %s
	fileStoragePath : %s
	key             : %s
	pollInterval    : %v
	reportInterval  : %v
	restore         : %v
	serverAddress   : %s
	storeInterval   : %v
	urlHost         : %s
`
	return fmt.Sprintf(format,
		c.DataBaseDSN(),
		c.FileStoragePath(),
		util.Str(c.Key()),
		c.PollInterval(),
		c.ReportInterval(),
		c.Restore(),
		c.ServerAddress(),
		c.StoreInterval(),
		util.Str(c.URLHost()),
	)
}

func (c *config) URLHost() *string {

	if c.urlHost != nil {
		return c.urlHost
	}
	//goland:noinspection HttpUrlsUsage
	urlHost := fmt.Sprintf("http://%s", c.serverAddress)
	c.urlHost = &urlHost

	return c.urlHost
}

func getConfig(opts ...func(*config)) *config {

	cfg = new(config)

	// вызываем все указанные функции для установки параметров
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func initAgentConfig() {

	initServerAddress()
	if env.ReportInterval < 1 {
		cfg.reportInterval = *flag.reportInterval
	} else {
		cfg.reportInterval = time.Duration(env.ReportInterval)
	}
	if env.PollInterval < 1 {
		cfg.pollInterval = *flag.pollInterval
	} else {
		cfg.pollInterval = time.Duration(env.PollInterval)
	}
	initKey()
}

func initServerConfig() {

	initServerAddress()
	if env.DataBaseDSN == "" {
		cfg.dataBaseDSN = flag.dataBaseDSN
	} else {
		cfg.dataBaseDSN = &env.DataBaseDSN
	}
	if env.StoreInterval == "" {
		cfg.storeInterval = *flag.storeInterval
	} else {
		storeInterval, err := strconv.Atoi(env.StoreInterval)
		if err == nil {
			cfg.storeInterval = time.Duration(storeInterval)
		} else {
			cfg.storeInterval = 24 * time.Hour
		}
	}
	if env.FileStoragePath == "" {
		cfg.fileStoragePath = *flag.fileStoragePath
	} else {
		cfg.fileStoragePath = env.FileStoragePath
	}
	if env.Restore == "" {
		cfg.restore = *flag.restore
	} else {
		restore, err := strconv.ParseBool(env.Restore)
		if err == nil {
			cfg.restore = restore
		}
	}
	initKey()
}

func initKey() {
	if env.Key == "" {
		cfg.key = flag.key
	} else {
		cfg.key = &env.Key
	}
}

func initServerAddress() {
	if len(env.Address) < 2 {
		cfg.serverAddress = *flag.serverAddress
	} else {
		cfg.serverAddress = env.parseEnvAddress()
	}
}
