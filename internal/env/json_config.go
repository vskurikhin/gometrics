/*
 * This file was last modified at 2024-06-24 16:57 by Victor N. Skurikhin.
 * json_config.go
 * $Id$
 */

package env

type jsConfig interface {
	getAddress() string
}

var jsonConfig jsConfig
