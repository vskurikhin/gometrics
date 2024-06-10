/*
 * This file was last modified at 2024-06-10 22:55 by Victor N. Skurikhin.
 * flags.go
 * $Id$
 */

package env

import (
	"github.com/spf13/pflag"
	"sync"
	"time"
)

type flags struct {
	serverAddress   *string
	storeInterval   *time.Duration
	fileStoragePath *string
	restore         *bool
	dataBaseDSN     *string
	urlHost         *string
	reportInterval  *time.Duration
	pollInterval    *time.Duration
	key             *string
}

var onceFlags = new(sync.Once)
var flag *flags

func initAgentFlags() {

	onceFlags.Do(func() {

		flag = new(flags)
		if !pflag.Parsed() {
			poll := pflag.IntP(
				"poll-interval",
				"p",
				2,
				"help message for poll interval",
			)
			report := pflag.IntP(
				"report-interval",
				"r",
				10,
				"help message for report interval",
			)
			initServerAddressFlag()
			initKeyFlag()
			pflag.Parse()

			pollInterval := time.Duration(*poll)
			flag.pollInterval = &pollInterval

			reportInterval := time.Duration(*report)
			flag.reportInterval = &reportInterval
		} else {
			pollInterval := 2 * time.Minute
			flag.pollInterval = &pollInterval

			reportInterval := 10 * time.Minute
			flag.reportInterval = &reportInterval
		}
	})
}

func initKeyFlag() {
	flag.key = pflag.StringP(
		"key",
		"k",
		"",
		"help message for key",
	)
}

func initServerAddressFlag() {
	flag.serverAddress = pflag.StringP(
		"address",
		"a",
		"localhost:8080",
		"help message for host and port",
	)
}

func initServerFlags() {

	onceFlags.Do(func() {

		flag = new(flags)

		if !pflag.Parsed() {
			sInterval := pflag.IntP(
				"store-interval",
				"i",
				300,
				"help message for store interval",
			)
			flag.fileStoragePath = pflag.StringP(
				"file-storage-path",
				"f",
				"/tmp/metrics-db.json",
				"help message for file storage path",
			)
			flag.restore = pflag.BoolP(
				"restore",
				"r",
				true,
				"help message for restore trigger",
			)
			flag.dataBaseDSN = pflag.StringP(
				"database-dsn",
				"d",
				"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable",
				"help message for file database DSN",
			)
			initServerAddressFlag()
			initKeyFlag()
			pflag.Parse()

			storeInterval := time.Duration(*sInterval)
			flag.storeInterval = &storeInterval
		}
	})
}
