/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * environments.go
 * $Id$
 */

package env

import (
	"fmt"
	"strconv"
	"sync"

	c0env "github.com/caarlos0/env"

	"github.com/vskurikhin/gometrics/internal/util"
)

type environments struct {
	Address         []string `env:"ADDRESS" envSeparator:":"`
	Config          string   `env:"CONFIG"`
	CryptoKey       string   `env:"CRYPTO_KEY" envSeparator:"-"`
	DNS             string   `env:"DNS"`
	DataBaseDSN     string   `env:"DATABASE_DSN"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH"`
	GRPCAddress     []string `env:"GRPC_ADDRESS" envSeparator:":"`
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

func (e *environments) parseEnvGRPCAddress() string {

	port, err := strconv.Atoi(e.GRPCAddress[1])
	util.IfErrorThenPanic(err)
	address := fmt.Sprintf("%s:%d", e.GRPCAddress[0], port)

	return address
}
