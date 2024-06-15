/*
 * This file was last modified at 2024-06-16 14:00 by Victor N. Skurikhin.
 * config.go
 * $Id$
 */

// Package env работа с настройками и окружением
package env

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/util"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config interface {
	fmt.Stringer
	CryptoKey() []string
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
	cryptoKey       []string
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

// GetAgentConfig — конфигурация для агента.
func GetAgentConfig() Config {
	onceCfg.Do(func() {
		cfg = new(config)
		getEnvironments()
		initAgentFlags()
		initAgentConfig()
	})
	return cfg
}

// GetServerConfig — конфигурация для сервера.
func GetServerConfig() Config {
	onceCfg.Do(func() {
		cfg = new(config)
		getEnvironments()
		initServerFlags()
		initServerConfig()
	})
	return cfg
}

// GetTestConfig — для создания тестовой конфигурации.
func GetTestConfig(opts ...func(*config)) Config {
	return getConfig(opts...)
}

// WithCryptoKey — поддержка асимметричного шифрования.
func WithCryptoKey(cryptoKey []string) func(*config) {
	return func(c *config) {
		c.cryptoKey = cryptoKey
	}
}

// CryptoKey — поддержка асимметричного шифрования.
func (c *config) CryptoKey() []string {
	return c.cryptoKey
}

// WithDataBaseDSN — строка для конфигурации подключения к БД.
func WithDataBaseDSN(dataBaseDSN *string) func(*config) {
	return func(c *config) {
		c.dataBaseDSN = dataBaseDSN
	}
}

// DataBaseDSN - геттер строки для конфигурации подключения к БД.
func (c *config) DataBaseDSN() string {
	if c.dataBaseDSN != nil {
		return *c.dataBaseDSN
	}
	return ""
}

// WithFileStoragePath — имя файла, куда сохраняются текущие значения.
// (по умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск).
func WithFileStoragePath(fileStoragePath string) func(*config) {
	return func(c *config) {
		c.fileStoragePath = fileStoragePath
	}
}

// FileStoragePath - геттер имени файла, куда сохраняются текущие значения.
func (c *config) FileStoragePath() string {
	return c.fileStoragePath
}

// IsDBSetup - геттер признака сконфигурированного соединения с БД.
func (c *config) IsDBSetup() bool {
	return c.dataBaseDSN != nil && *c.dataBaseDSN != ""
}

// WithKey — ключ для вычисления и передачи хеша в HTTP-заголовке запроса с именем HashSHA25.
func WithKey(key *string) func(*config) {
	return func(c *config) {
		c.key = key
	}
}

// Key - геттер ключа вычисления и передачи хеша в HTTP-заголовке запроса с именем HashSHA25.
func (c *config) Key() *string {
	return c.key
}

// WithPollInterval — частота опроса метрик из пакета runtime (по умолчанию 2 секунды).
func WithPollInterval(pollInterval time.Duration) func(*config) {
	return func(c *config) {
		c.pollInterval = pollInterval
	}
}

// PollInterval — геттер для частоты  опроса метрик из пакета runtime.
func (c *config) PollInterval() time.Duration {
	return c.pollInterval
}

// WithReportInterval — частота отправки метрик на сервер (по умолчанию 10 секунд).
func WithReportInterval(reportInterval time.Duration) func(*config) {
	return func(c *config) {
		c.reportInterval = reportInterval
	}
}

// ReportInterval — геттер для частоты отправки метрик на сервер.
func (c *config) ReportInterval() time.Duration {
	return c.reportInterval
}

// WithRestore — булев признак (true/false), определяющей, загружать или нет ранее сохранённые
// значения из указанного файла при старте сервера (по умолчанию true).
func WithRestore(restore bool) func(*config) {
	return func(c *config) {
		c.restore = restore
	}
}

// Restore - геттер признака для загрузки сохранённых значений из файла.
func (c *config) Restore() bool {
	return c.restore
}

// WithServerAddress — адрес эндпоинта HTTP-сервера (по умолчанию localhost:8080).
func WithServerAddress(serverAddress string) func(*config) {
	return func(c *config) {
		c.serverAddress = serverAddress
	}
}

// ServerAddress - геттер для адреса HTTP-сервера.
func (c *config) ServerAddress() string {
	return c.serverAddress
}

// WithStoreInterval — интервал времени в секундах, по истечении которого текущие показания
// сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной).
func WithStoreInterval(storeInterval time.Duration) func(*config) {
	return func(c *config) {
		c.storeInterval = storeInterval
	}
}

// StoreInterval — геттер интервала сервера для сохранения на диск.
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

	initAgentCryptoKey()
	initServerAddress()
	if env.ReportInterval < 1 {
		cfg.reportInterval = *flag.reportInterval
	} else {
		cfg.reportInterval = time.Duration(env.ReportInterval) * time.Second
	}
	if env.PollInterval < 1 {
		cfg.pollInterval = *flag.pollInterval
	} else {
		cfg.pollInterval = time.Duration(env.PollInterval) * time.Second
	}
	initKey()
}

func initServerConfig() {

	initServerAddress()
	initServerCryptoKey()
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
			cfg.storeInterval = time.Duration(storeInterval) * time.Second
		} else {
			cfg.storeInterval = 300 * time.Second
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

func initAgentCryptoKey() {
	cfg.cryptoKey = []string{*flag.cryptoKey, "010001"}
}

func initServerCryptoKey() {
	cfg.cryptoKey = strings.Split(*flag.cryptoKey, "-")
}

func initServerAddress() {
	if len(env.Address) < 2 {
		cfg.serverAddress = *flag.serverAddress
	} else {
		cfg.serverAddress = env.parseEnvAddress()
	}
}
