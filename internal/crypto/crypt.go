/*
 * This file was last modified at 2024-07-05 16:00 by Victor N. Skurikhin.
 * crypt.go
 * $Id$
 */

// Package crypto поддержка асимметричного шифрования.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/util"
	"sync"
)

var _ Crypto = (*crypto)(nil)
var once = new(sync.Once)
var crypt *crypto

type crypto struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type Crypto interface {
	EncryptRSA(plain []byte) ([]byte, error)
	DecryptRSA(bytes []byte) ([]byte, error)
	EncryptAES(plain []byte) ([]byte, []byte, error)
	DecryptAES(secretKey, bytes []byte) ([]byte, error)
}

func GetAgentCrypto() Crypto {
	once.Do(func() {
		crypt = new(crypto)
		parameters := env.GetParameters()
		crypt.publicKey = parameters.PublicKey()
	})
	return crypt
}

func GetServerCrypto() Crypto {
	once.Do(func() {
		crypt = new(crypto)
		parameters := env.GetParameters()
		crypt.privateKey = parameters.PrivateKey()
	})
	return crypt
}

func (c *crypto) EncryptRSA(plain []byte) ([]byte, error) {

	if c.publicKey != nil {

		if result, err := rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, plain); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
	return nil, fmt.Errorf("public key is null")
}

func (c *crypto) DecryptRSA(bytes []byte) ([]byte, error) {

	if c.privateKey != nil {

		if result, err := rsa.DecryptPKCS1v15(nil, c.privateKey, bytes); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
	return nil, fmt.Errorf("private key is null")
}

func (c *crypto) EncryptAES(plain []byte) ([]byte, []byte, error) {

	if c.publicKey != nil {

		var gcm cipher.AEAD
		secretKey := make([]byte, 32) // 32 bytes to select AES-256.

		if _, err := rand.Reader.Read(secretKey); err != nil {
			return nil, nil, err
		}
		if block, err := aes.NewCipher(secretKey); err != nil {
			return nil, nil, err
		} else {
			gcm, err = cipher.NewGCM(block)
			if err != nil {
				return nil, nil, err
			}
		}
		// Нужен 12-байтовый nonce для GCM (можно изменить, если вы используете cipher.NewGCMWithNonceSize()).
		// nonce всегда должен генерироваться случайным образом для каждого шифрования.
		nonce := make([]byte, gcm.NonceSize())

		if _, err := rand.Read(nonce); err != nil {
			return nil, nil, err
		}
		// Зашифрованные данные на самом деле nonce+зашифрованный plain.
		// При расшифровке просто зная размер nonce,
		// достаточно, чтобы отделить его от зашифрованных данных.
		return secretKey, gcm.Seal(nonce, nonce, plain, nil), nil
	}
	return nil, nil, fmt.Errorf("public key is null")
}

func (c *crypto) DecryptAES(secretKey, bytes []byte) ([]byte, error) {

	var gcm cipher.AEAD

	if block, err := aes.NewCipher(secretKey); err != nil {
		return nil, err
	} else {
		if gcm, err = cipher.NewGCM(block); err != nil {
			return nil, err
		}
	}
	// Так как зашифрованные данные представляют собой nonce + зашифрованный plain
	// и len(nonce) == NonceSize(). Можно их разделить.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := bytes[:nonceSize], bytes[nonceSize:]

	if plain, err := gcm.Open(nil, nonce, ciphertext, nil); err != nil {
		return nil, err
	} else {
		return plain, nil
	}
}

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	util.IfErrorThenPanic(err)
	return privateKey, &privateKey.PublicKey
}
