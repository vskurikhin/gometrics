/*
 * This file was last modified at 2024-07-02 15:12 by Victor N. Skurikhin.
 * init_config_test.go
 * $Id$
 */

package env

import (
	c0env "github.com/caarlos0/env"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitAgentCryptoKey(t *testing.T) {
	t.Setenv("CRYPTO_KEY", testCryptoKey)
	cfg = new(config)
	env = new(environments)
	err := c0env.Parse(env)
	initAgentCryptoKey()
	assert.Nil(t, err)
	assert.Equal(t, testCryptoKey, cfg.CryptoKey())
}

func TestInitServerCryptoKey(t *testing.T) {
	t.Setenv("CRYPTO_KEY", testCryptoKey)
	cfg = new(config)
	env = new(environments)
	err := c0env.Parse(env)
	initServerCryptoKey()
	assert.Nil(t, err)
	assert.Equal(t, testCryptoKey, cfg.CryptoKey())
}
