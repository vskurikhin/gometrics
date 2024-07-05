/*
 * This file was last modified at 2024-07-04 17:29 by Victor N. Skurikhin.
 * environments.go
 * $Id$
 */

package env

import (
	"fmt"
	c0env "github.com/caarlos0/env"
	"github.com/vskurikhin/gometrics/internal/util"
	"strconv"
	"sync"
)

type environments struct {
	Address         []string `env:"ADDRESS" envSeparator:":"`
	Config          string   `env:"CONFIG"`
	CryptoKey       string   `env:"CRYPTO_KEY" envSeparator:"-"`
	DataBaseDSN     string   `env:"DATABASE_DSN"`
	DNS             string   `env:"DNS"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH"`
	Key             string   `env:"KEY"`
	PollInterval    int      `env:"POLL_INTERVAL"`
	ReportInterval  int      `env:"REPORT_INTERVAL"`
	Restore         string   `env:"RESTORE"`
	StoreInterval   string   `env:"STORE_INTERVAL"`
	TrustedSubnet   string   `env:"TRUSTED_SUBNET"`
}

var onceEnv = new(sync.Once)
var env *environments

func getEnvironments() *environments {

	onceEnv.Do(func() {
		env = new(environments)
		err := c0env.Parse(env)
		util.IfErrorThenPanic(err)
		if env.DNS == "" {
			env.DNS = "8.8.8.8:53"
		}
	})
	return env
}

func (e *environments) parseEnvAddress() string {

	port, err := strconv.Atoi(e.Address[1])
	util.IfErrorThenPanic(err)
	address := fmt.Sprintf("%s:%d", e.Address[0], port)

	return address
}
