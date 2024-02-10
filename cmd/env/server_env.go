/*
 * This file was last modified at 2024-02-10 23:44 by Victor N. Skurikhin.
 * server_env.go
 * $Id$
 */

package env

type serverEnv struct {
	serverAddress *string
}

func (sf *serverEnv) ServerAddress() string {
	return *sf.serverAddress
}
