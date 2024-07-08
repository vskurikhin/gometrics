/*
 * This file was last modified at 2024-07-07 11:44 by Victor N. Skurikhin.
 * json_config.go
 * $Id$
 */

package env

type jsConfig interface {
	getAddress() string
	getGRPCAddress() string
}

var jsonConfig jsConfig
