/*
 * This file was last modified at 2024-06-24 16:57 by Victor N. Skurikhin.
 * json_agent_config.go
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

// agentConfig конфигурации агента с помощью файла в формате JSON.
type agentConfig struct {
	Address        string `json:"address"`         // Address аналог переменной окружения ADDRESS или флага -a
	ReportInterval string `json:"report_interval"` // ReportInterval аналог переменной окружения REPORT_INTERVAL или флага -r
	PollInterval   string `json:"poll_interval"`   // PollInterval аналог переменной окружения POLL_INTERVAL или флага -p
	CryptoKey      string `json:"crypto_key"`      // CryptoKey аналог переменной окружения CRYPTO_KEY или флага -crypto-key
}

var jsonAgentConfig *agentConfig

func getAgentConfig() *agentConfig {

	jsonAgentConfig = new(agentConfig)

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

		if err := easyjson.UnmarshalFromReader(buf, jsonAgentConfig); err != nil {
			panic(err)
		}
	}
	jsonConfig = jsonAgentConfig

	return jsonAgentConfig
}

func (a *agentConfig) getAddress() string {
	return a.Address
}
