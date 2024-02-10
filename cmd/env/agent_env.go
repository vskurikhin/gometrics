/*
 * This file was last modified at 2024-02-10 23:44 by Victor N. Skurikhin.
 * agent_env.go
 * $Id$
 */

package env

import (
	"fmt"
	"time"
)

type agentEnv struct {
	serverEnv
	urlHost        *string
	reportInterval time.Duration
	pollInterval   time.Duration
}

func (ae *agentEnv) URLHost() *string {

	if ae.urlHost != nil {
		return ae.urlHost
	}
	//goland:noinspection HttpUrlsUsage
	urlHost := fmt.Sprintf("http://%s", *ae.serverAddress)
	ae.urlHost = &urlHost

	return ae.urlHost
}

func (ae *agentEnv) ReportInterval() time.Duration {
	return ae.reportInterval
}

func (ae *agentEnv) PollInterval() time.Duration {
	return ae.pollInterval
}
