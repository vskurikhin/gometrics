/*
 * This file was last modified at 2024-03-18 11:11 by Victor N. Skurikhin.
 * init.go
 * $Id$
 */

package env

import (
	"fmt"
	"github.com/caarlos0/env"
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
		Server.dataBaseDSN = *sFlags.dataBaseDSN
	} else {
		Server.dataBaseDSN = cfg.DataBaseDSN
	}
}
