/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * consts.go
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
