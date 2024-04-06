/*
 * This file was last modified at 2024-04-06 18:15 by Victor N. Skurikhin.
 * init.go
 * $Id$
 */

package env

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"
)

var (
	aFlags = agentFlags{}
	sFlags = serverFlags{}
	Agent  = agentEnv{}
	Server = serverEnv{}
)

var cfg config

func InitAgent() {

	initAgentFlags()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.Address) < 2 {
		Agent.urlHost = aFlags.URLHost()
	} else {
		address := parseEnvAddress(cfg)
		Agent.serverAddress = &address
	}

	if cfg.ReportInterval < 1 {
		Agent.reportInterval = aFlags.ReportInterval()
	} else {
		Agent.reportInterval = time.Duration(cfg.ReportInterval)
	}

	if cfg.PollInterval < 1 {
		Agent.pollInterval = aFlags.PollInterval()
	} else {
		Agent.pollInterval = time.Duration(cfg.PollInterval)
	}

	if cfg.Key == "" {
		Agent.key = aFlags.key
	} else {
		Agent.key = &cfg.Key
	}

	if cfg.RateLimit < 1 && aFlags.RateLimit() != nil {
		Agent.rateLimit = *aFlags.RateLimit()
	} else {
		Agent.rateLimit = cfg.RateLimit
	}
	logger.Log.Debug("Agent ", zap.String("env", fmt.Sprintf("%+v", Agent)))
}

func parseEnvAddress(cfg config) string {

	port, err := strconv.Atoi(cfg.Address[1])
	if err != nil {
		panic(err)
	}
	address := fmt.Sprintf("%s:%d", cfg.Address[0], port)

	return address
}

func InitServer() {

	initServerFlags()

	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if len(cfg.Address) < 2 {
		address := sFlags.ServerAddress()
		Server.serverAddress = &address
	} else {
		address := parseEnvAddress(cfg)
		Server.serverAddress = &address
	}

	if cfg.StoreInterval == "" {
		Server.storeInterval = sFlags.StoreInterval()
	} else {
		storeInterval, err := strconv.Atoi(cfg.StoreInterval)
		if err == nil {
			Server.storeInterval = time.Duration(storeInterval)
		}
	}

	if cfg.FileStoragePath == "" {
		Server.fileStoragePath = *sFlags.fileStoragePath
	} else {
		Server.fileStoragePath = cfg.FileStoragePath
	}

	if cfg.Restore == "" {
		Server.restore = *sFlags.restore
	} else {
		restore, err := strconv.ParseBool(cfg.Restore)
		if err == nil {
			Server.restore = restore
		}
	}

	if cfg.DataBaseDSN == "" {
		Server.dataBaseDSN = sFlags.dataBaseDSN
	} else {
		Server.dataBaseDSN = &cfg.DataBaseDSN
	}

	if cfg.Key == "" {
		Server.key = sFlags.key
	} else {
		Server.key = &cfg.Key
	}
	logger.Log.Debug("Server ", zap.String("env", fmt.Sprintf("%+v", Server)))
}
