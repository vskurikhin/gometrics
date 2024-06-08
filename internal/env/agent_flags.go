/*
 * This file was last modified at 2024-05-28 16:19 by Victor N. Skurikhin.
 * agent_flags.go
 * $Id$
 */

package env

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

type agentFlags struct {
	serverFlags
	urlHost        *string
	reportInterval *time.Duration
	pollInterval   *time.Duration
	key            *string
}

func (af *agentFlags) URLHost() *string {

	if af.urlHost != nil {
		return af.urlHost
	}
	//goland:noinspection HttpUrlsUsage
	urlHost := fmt.Sprintf("http://%s", *af.serverAddress)
	af.urlHost = &urlHost

	return af.urlHost
}

func (af *agentFlags) ReportInterval() time.Duration {
	return *af.reportInterval
}

func (af *agentFlags) PollInterval() time.Duration {
	return *af.pollInterval
}

func initAgentFlags() {

	aFlags.serverAddress = pflag.StringP("address", "a", "localhost:8080", "help message for host and port")

	report := pflag.IntP("report-interval", "r", 10, "help message for report interval")
	poll := pflag.IntP("poll-interval", "p", 2, "help message for poll interval")
	aFlags.key = pflag.StringP("key", "k", "", "help message for key")

	pflag.Parse()

	reportInterval := time.Duration(*report)
	aFlags.reportInterval = &reportInterval

	pollInterval := time.Duration(*poll)
	aFlags.pollInterval = &pollInterval

}
