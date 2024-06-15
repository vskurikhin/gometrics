/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
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
	ReportInterval  int      `env:"REPORT_INTERVAL"`
	PollInterval    int      `env:"POLL_INTERVAL"`
	StoreInterval   string   `env:"STORE_INTERVAL"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH"`
	Restore         string   `env:"RESTORE"`
	DataBaseDSN     string   `env:"DATABASE_DSN"`
	Key             string   `env:"KEY"`
}

var onceEnv = new(sync.Once)
var env *environments

func getEnvironments() *environments {

	onceEnv.Do(func() {
		env = new(environments)
		err := c0env.Parse(env)
		util.IfErrorThenPanic(err)
	})
	return env
}

func (e *environments) parseEnvAddress() string {

	port, err := strconv.Atoi(e.Address[1])
	util.IfErrorThenPanic(err)
	address := fmt.Sprintf("%s:%d", e.Address[0], port)

	return address
}
