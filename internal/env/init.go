/*
 * This file was last modified at 2024-02-10 23:59 by Victor N. Skurikhin.
 * init.go
 * $Id$
 */

package env

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
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

func InitAgent() {

	initAgentFlags()

	var cfg config
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
}

func initAgentFlags() {

	aFlags.serverAddress = pflag.StringP("address", "a", "localhost:8080", "help message for host and port")

	report := pflag.IntP("report-interval", "r", 10, "help message for report interval")
	poll := pflag.IntP("poll-interval", "p", 2, "help message for poll interval")

	pflag.Parse()

	reportInterval := time.Duration(*report)
	aFlags.reportInterval = &reportInterval

	pollInterval := time.Duration(*poll)
	aFlags.pollInterval = &pollInterval
}

func initServerFlags() {

	sFlags.serverAddress = pflag.StringP("address", "a", "localhost:8080", "help message for host and port")
	pflag.Parse()
}
