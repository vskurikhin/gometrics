/*
 * This file was last modified at 2024-02-04 13:32 by Victor N. Skurikhin.
 * api_name.go
 * $Id$
 */

package names

const (
	Port            = 8080
	Host            = "http://localhost:8080"
	Value           = "/value"
	ValueURL        = Value + "/"
	ValueChi        = ValueURL + "{type}/{name}"
	Update          = "/update"
	UpdateURL       = Update + "/"
	UpdateChi       = UpdateURL + "{type}/{name}/{value:[a-zA-Z0-9-+.]+}"
	UpdateURLClient = Host + Update
)
