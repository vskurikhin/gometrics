/*
 * This file was last modified at 2024-03-18 22:06 by Victor N. Skurikhin.
 * server_flags.go
 * $Id$
 */

package env

import (
	"github.com/spf13/pflag"
	"time"
)

type serverFlags struct {
	serverAddress   *string
	storeInterval   *time.Duration
	fileStoragePath *string
	restore         *bool
	dataBaseDSN     *string
	key             *string
}

func (sf *serverFlags) ServerAddress() string {
	return *sf.serverAddress
}

func (sf *serverFlags) StoreInterval() time.Duration {
	return *sf.storeInterval
}

func (sf *serverFlags) FileStoragePath() string {
	return *sf.fileStoragePath
}

func (sf *serverFlags) Restore() bool {
	return *sf.restore
}

func (sf *serverFlags) DataBaseDSN() string {
	return *sf.dataBaseDSN
}

func initServerFlags() {

	sFlags.serverAddress = pflag.StringP("address", "a", "localhost:8080", "help message for host and port")
	sInterval := pflag.IntP("store-interval", "i", 300, "help message for store interval")
	sFlags.fileStoragePath = pflag.StringP("file-storage-path", "f", "/tmp/metrics-db.json", "help message for file storage path")
	sFlags.restore = pflag.BoolP("restore", "r", true, "help message for restore trigger")
	sFlags.dataBaseDSN = pflag.StringP("database-dsn", "d", "", "help message for file database DSN")
	sFlags.key = pflag.StringP("key", "k", "", "help message for key")

	pflag.Parse()

	storeInterval := time.Duration(*sInterval)
	sFlags.storeInterval = &storeInterval

}
