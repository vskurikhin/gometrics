/*
 * This file was last modified at 2024-02-10 16:09 by Victor N. Skurikhin.
 * agent.go
 * $Id$
 */

package cflag

import (
	"fmt"
	"time"
)

type agentFlags struct {
	serverFlags
	urlHost        *string
	reportInterval *int
	pollInterval   *int
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
	return time.Duration(*af.reportInterval)
}

func (af *agentFlags) PollInterval() time.Duration {
	return time.Duration(*af.pollInterval)
}
