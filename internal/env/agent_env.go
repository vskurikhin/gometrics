/*
 * This file was last modified at 2024-04-06 16:54 by Victor N. Skurikhin.
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
	key            *string
	rateLimit      int
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

func (ae *agentEnv) Key() *string {
	return ae.key
}

func (ae *agentEnv) RateLimit() int {
	return ae.rateLimit
}
