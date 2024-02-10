/*
 * This file was last modified at 2024-02-10 23:44 by Victor N. Skurikhin.
 * config.go
 * $Id$
 */

package env

type config struct {
	Address        []string `env:"ADDRESS" envSeparator:":"`
	ReportInterval int      `env:"REPORT_INTERVAL"`
	PollInterval   int      `env:"POLL_INTERVAL"`
}
