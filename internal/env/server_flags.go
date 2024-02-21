/*
 * This file was last modified at 2024-02-10 23:59 by Victor N. Skurikhin.
 * server_flags.go
 * $Id$
 */

package env

type serverFlags struct {
	serverAddress *string
}

func (sf *serverFlags) ServerAddress() string {
	return *sf.serverAddress
}
