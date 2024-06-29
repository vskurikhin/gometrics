/*
 * This file was last modified at 2024-06-24 16:57 by Victor N. Skurikhin.
 * json_server_config.go
 * $Id$
 */

package env

import (
	"bytes"
	"github.com/mailru/easyjson"
	"github.com/vskurikhin/gometrics/internal/util"
	"io"
	"os"
)

// serverConfig конфигурация сервера с помощью файла в формате JSON.
type serverConfig struct {
	Address       string `json:"address"`        // Address аналог переменной окружения ADDRESS или флага -a
	Restore       bool   `json:"restore"`        // Restore аналог переменной окружения RESTORE или флага -r
	StoreInterval string `json:"store_interval"` // StoreInterval аналог переменной окружения STORE_INTERVAL или флага -i
	StoreFile     string `json:"store_file"`     // StoreFile аналог переменной окружения STORE_FILE или -f
	DatabaseDSN   string `json:"database_dsn"`   // DatabaseDSN аналог переменной окружения DATABASE_DSN или флага -d
	CryptoKey     string `json:"crypto_key"`     // CryptoKey аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}

var jsonServerConfig *serverConfig

func getServerConfig() *serverConfig {

	jsonServerConfig = new(serverConfig)

	if flag.config != nil && *flag.config != "" {
		cfg.configFileName = *flag.config
	} else if env.Config != "" {
		cfg.configFileName = env.Config
	}
	if cfg.configFileName != "" {

		configFile, err := os.Open(cfg.configFileName)
		util.IfErrorThenPanic(err)
		b, err := io.ReadAll(configFile)
		util.IfErrorThenPanic(err)
		buf := bytes.NewBuffer(b)

		if err := easyjson.UnmarshalFromReader(buf, jsonServerConfig); err != nil {
			panic(err)
		}
	}
	jsonConfig = jsonServerConfig

	return jsonServerConfig
}

func (sc *serverConfig) getAddress() string {
	return sc.Address
}
