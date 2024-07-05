/*
 * This file was last modified at 2024-07-03 20:58 by Victor N. Skurikhin.
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
	config          *string
	cryptoKey       *string
	dataBaseDSN     *string
	fileStoragePath *string
	key             *string
	pollInterval    *time.Duration
	reportInterval  *time.Duration
	restore         *bool
	serverAddress   *string
	storeInterval   *time.Duration
	trustedSubnet   *string
	urlHost         *string
}

var onceFlags = new(sync.Once)
var flag *flags

func initAgentFlags() {

	onceFlags.Do(func() {

		flag = new(flags)
		if !pflag.Parsed() {
			poll := pflag.IntP(
				PollInterval,
				"p",
				2,
				"help message for poll interval",
			)
			report := pflag.IntP(
				ReportInterval,
				"r",
				10,
				"help message for report interval",
			)
			initConfigFlag()
			initCryptoAgentFlag()
			initServerAddressFlag()
			initKeyFlag()
			pflag.Parse()

			pollInterval := time.Duration(*poll) * time.Second
			flag.pollInterval = &pollInterval

			reportInterval := time.Duration(*report) * time.Second
			flag.reportInterval = &reportInterval
		} else {
			pollInterval := 2 * time.Minute
			flag.pollInterval = &pollInterval

			reportInterval := 10 * time.Minute
			flag.reportInterval = &reportInterval
		}
	})
}

func initConfigFlag() {
	flag.config = pflag.StringP(
		ConfigFileName,
		"c",
		"",
		"help message for Config file",
	)
}

func initCryptoAgentFlag() {
	flag.cryptoKey = pflag.StringP(
		CryptoKey,
		"y",
		"public_key.pem",
		"help message for crypto key",
	)
}

func initCryptoServerFlag() {
	flag.cryptoKey = pflag.StringP(
		CryptoKey,
		"y",
		"private_key.pem",
		"help message for crypto key",
	)
}

func initKeyFlag() {
	flag.key = pflag.StringP(
		Key,
		"k",
		"",
		"help message for key",
	)
}

func initServerAddressFlag() {
	flag.serverAddress = pflag.StringP(
		Address,
		"a",
		"localhost:8080",
		"help message for host and port",
	)
}

func initServerFlags() {

	onceFlags.Do(func() {

		flag = new(flags)

		if !pflag.Parsed() {
			flag.dataBaseDSN = pflag.StringP(
				DatabaseDSN,
				"d",
				"postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable",
				"help message for file database DSN",
			)
			sInterval := pflag.IntP(
				StoreInterval,
				"i",
				300,
				"help message for store interval",
			)
			flag.fileStoragePath = pflag.StringP(
				FileStoragePath,
				"f",
				"/tmp/metrics-db.json",
				"help message for file storage path",
			)
			flag.restore = pflag.BoolP(
				Restore,
				"r",
				true,
				"help message for restore trigger",
			)
			flag.trustedSubnet = pflag.StringP(
				TrustedSubnet,
				"t",
				"",
				"help message for trusted subnet",
			)
			initConfigFlag()
			initCryptoServerFlag()
			initServerAddressFlag()
			initKeyFlag()
			pflag.Parse()

			storeInterval := time.Duration(*sInterval) * time.Second
			flag.storeInterval = &storeInterval
		}
	})
}

func setIfFlagChanged(name string, set func()) {
	pflag.VisitAll(func(f *pflag.Flag) {
		if f.Changed && f.Name == name {
			set()
		}
		if f.Name == ConfigFileName && !f.Changed && env.Config == "" {
			set()
		}
	})
}
