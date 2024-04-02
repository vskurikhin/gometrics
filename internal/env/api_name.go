/*
 * This file was last modified at 2024-04-03 08:47 by Victor N. Skurikhin.
 * api_name.go
 * $Id$
 */

package env

const (
	Port            = 8080
	Host            = "http://localhost:8080"
	Ping            = "/ping"
	PingURL         = Ping + "/"
	Value           = "/value"
	ValueURL        = Value + "/"
	ValueChi        = ValueURL + "{type}/{name}"
	Update          = "/update"
	UpdateURL       = Update + "/"
	Updates         = "/updates"
	UpdatesURL      = Updates + "/"
	UpdateChi       = UpdateURL + "{type}/{name}/{value:[a-zA-Z0-9-+.]+}"
	UpdateURLClient = Host + Update
)
