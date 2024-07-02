/*
 * This file was last modified at 2024-07-02 15:12 by Victor N. Skurikhin.
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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/util"
	"io"
	"os"
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

func GetAgentCrypto(cfg env.Config) Crypto {
	once.Do(func() {
		crypt = new(crypto)
		if len(cfg.CryptoKey()) > 1 {
			file, err := os.Open(cfg.CryptoKey())
			if err != nil {
				return
			}
			//nolint:multichecker,errcheck
			defer func() { _ = file.Close() }()
			buf, err := io.ReadAll(file)
			if err != nil {
				return
			}
			if block := readPEMString(string(buf)); block != nil {
				publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
				util.IfErrorThenPanic(err)
				crypt.publicKey = publicKey
			}
		}
	})
	return crypt
}

func GetServerCrypto(cfg env.Config) Crypto {
	once.Do(func() {
		crypt = new(crypto)
		if len(cfg.CryptoKey()) > 1 {
			file, err := os.Open(cfg.CryptoKey())
			if err != nil {
				return
			}
			//nolint:multichecker,errcheck
			defer func() { _ = file.Close() }()
			buf, err := io.ReadAll(file)
			if err != nil {
				return
			}
			if block := readPEMString(string(buf)); block != nil {
				privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
				util.IfErrorThenPanic(err)
				crypt.privateKey = privateKey
			}
		}
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
		// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
		// A nonce should always be randomly generated for every encryption.
		nonce := make([]byte, gcm.NonceSize())

		if _, err := rand.Read(nonce); err != nil {
			return nil, nil, err
		}
		// ciphertext here is actually nonce+ciphertext
		// So that when we decrypt, just knowing the nonce size
		// is enough to separate it from the ciphertext.
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
	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := bytes[:nonceSize], bytes[nonceSize:]

	if plain, err := gcm.Open(nil, nonce, ciphertext, nil); err != nil {
		return nil, err
	} else {
		return plain, nil
	}
}

func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	util.IfErrorThenPanic(err)
	return privateKey, &privateKey.PublicKey
}

func readPEMString(p string) *pem.Block {
	result, _ := pem.Decode([]byte(p))
	return result
}
