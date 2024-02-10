/*
 * This file was last modified at 2024-02-10 16:23 by Victor N. Skurikhin.
 * init.go
 * $Id$
 */

package cflag

import (
	"github.com/spf13/pflag"
)

var (
	AgentFlags  = agentFlags{}
	ServerFlags = serverFlags{}
)

func init() {
	ServerFlags.serverAddress = pflag.StringP("address", "a", "localhost:8080", "help message for host and port")

}

func InitAgent() {
	AgentFlags.serverAddress = ServerFlags.serverAddress
	AgentFlags.reportInterval = pflag.IntP("report-interval", "r", 10, "help message for report interval")
	AgentFlags.pollInterval = pflag.IntP("poll-interval", "p", 2, "help message for poll interval")
}
