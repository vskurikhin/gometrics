/*
 * This file was last modified at 2024-02-10 23:59 by Victor N. Skurikhin.
 * agent_flags.go
 * $Id$
 */

package env

import (
	"fmt"
	"time"
)

type agentFlags struct {
	serverFlags
	urlHost        *string
	reportInterval *time.Duration
	pollInterval   *time.Duration
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
