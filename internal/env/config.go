/*
 * This file was last modified at 2024-03-02 13:51 by Victor N. Skurikhin.
 * config.go
 * $Id$
 */

package env

type config struct {
	Address         []string `env:"ADDRESS" envSeparator:":"`
	ReportInterval  int      `env:"REPORT_INTERVAL"`
	PollInterval    int      `env:"POLL_INTERVAL"`
	StoreInterval   string   `env:"STORE_INTERVAL"`
	FileStoragePath string   `env:"FILE_STORAGE_PATH"`
	Restore         string   `env:"RESTORE"`
}
