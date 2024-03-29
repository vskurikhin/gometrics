/*
 * This file was last modified at 2024-03-02 14:18 by Victor N. Skurikhin.
 * server_env.go
 * $Id$
 */

package env

import "time"

type serverEnv struct {
	serverAddress   *string
	storeInterval   time.Duration
	fileStoragePath string
	restore         bool
}

func (sf *serverEnv) ServerAddress() string {
	return *sf.serverAddress
}

func (sf *serverEnv) StoreInterval() time.Duration {
	return sf.storeInterval
}

func (sf *serverEnv) FileStoragePath() string {
	return sf.fileStoragePath
}

func (sf *serverEnv) Restore() bool {
	return sf.restore
}
