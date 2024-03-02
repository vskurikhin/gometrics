/*
 * This file was last modified at 2024-03-02 20:04 by Victor N. Skurikhin.
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

func initServerFlags() {

	sFlags.serverAddress = pflag.StringP("address", "a", "localhost:8080", "help message for host and port")
	sInterval := pflag.IntP("store-interval", "i", 300, "help message for store interval")
	sFlags.fileStoragePath = pflag.StringP("file-storage-path", "f", "/tmp/metrics-db.json", "help message for file storage path")
	sFlags.restore = pflag.BoolP("restore", "r", true, "help message for restore trigger")

	pflag.Parse()

	storeInterval := time.Duration(*sInterval)
	sFlags.storeInterval = &storeInterval
}
