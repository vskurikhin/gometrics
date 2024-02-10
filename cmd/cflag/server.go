/*
 * This file was last modified at 2024-02-10 15:07 by Victor N. Skurikhin.
 * server.go
 * $Id$
 */

package cflag

type serverFlags struct {
	serverAddress *string
}

func (sf *serverFlags) ServerAddress() string {
	return *sf.serverAddress
}
