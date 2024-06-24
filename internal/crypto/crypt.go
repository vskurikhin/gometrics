/*
 * This file was last modified at 2024-06-24 22:55 by Victor N. Skurikhin.
 * crypt.go
 * $Id$
 */

// Package crypto поддержка асимметричного шифрования.
package crypto

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/simplersa"
	"math/big"
	"sync"
)

var once = new(sync.Once)
var crypt *crypto

type crypto struct {
	privateKey *simplersa.PrivateKey
	publicKey  *simplersa.PublicKey
}

type Crypto interface {
	EncryptRSA(m []byte) ([]byte, error)
	TryDecryptRSA(b []byte) ([]byte, bool)
}

func GetAgentCrypto(cfg env.Config) Crypto {
	once.Do(func() {
		crypt = new(crypto)
		if len(cfg.CryptoKey()) > 1 {
			e := new(big.Int)
			n := new(big.Int)
			e.SetString(cfg.CryptoKey()[1], 16)
			n.SetString(cfg.CryptoKey()[0], 16)
			crypt.publicKey = &simplersa.PublicKey{N: n, E: e}
		}
	})
	return crypt
}

func GetServerCrypto(cfg env.Config) Crypto {
	once.Do(func() {
		crypt = new(crypto)
		if len(cfg.CryptoKey()) > 1 {
			d := new(big.Int)
			n := new(big.Int)
			d.SetString(cfg.CryptoKey()[1], 16)
			n.SetString(cfg.CryptoKey()[0], 16)
			crypt.privateKey = &simplersa.PrivateKey{N: n, D: d}
		}
	})
	return crypt
}

func (c *crypto) EncryptRSA(m []byte) ([]byte, error) {

	if c.publicKey != nil {
		if buf, err := simplersa.EncryptRSA(c.publicKey, m); err != nil {
			return nil, err
		} else {
			return buf, nil
		}
	}
	return nil, fmt.Errorf("public key is null")
}

func (c *crypto) TryDecryptRSA(b []byte) ([]byte, bool) {

	if c.privateKey != nil {
		if buf, err := simplersa.DecryptRSA(c.privateKey, b); err != nil {
			return nil, false
		} else {
			return buf, true
		}
	}
	return nil, false
}
