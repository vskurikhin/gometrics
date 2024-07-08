/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * crypt_test.go
 * $Id$
 */

package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/util"
)

var (
	expected               = "very long secret string"
	testPublicKeyFileName  string
	testPrivateKeyFileName string
)

func getTestConfigAgentCrypto() env.Config {
	return env.GetTestConfig(
		env.GetProperty,
		env.WithCryptoKey(testPublicKeyFileName),
	)
}

func getTestConfigServerCrypto() env.Config {
	return env.GetTestConfig(
		env.GetProperty,
		env.WithCryptoKey(testPrivateKeyFileName),
	)
}

func exportRsaPrivateKeyAsPemStr(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	return string(privateKeyPEM)
}

func exportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) (string, error) {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(pubkey)
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	return string(publicKeyPEM), nil
}

func TestEncryptRSA(t *testing.T) {
	getTestConfigAgentCrypto()
	t.Setenv("CRYPTO_KEY", testPublicKeyFileName)
	ca := GetAgentCrypto()
	cs := getTestServerCrypto()
	be, err := ca.EncryptRSA([]byte(expected))
	assert.Nil(t, err)
	assert.NotNil(t, be)
	got, e := cs.DecryptRSA(be)
	assert.Nil(t, e)
	assert.Equal(t, expected, string(got))
}

func TestEncryptAES(t *testing.T) {
	getTestConfigAgentCrypto()
	t.Setenv("CRYPTO_KEY", testPublicKeyFileName)
	ca := GetAgentCrypto()
	cs := getTestServerCrypto()
	secretKey, be, err := ca.EncryptAES([]byte(expected))
	assert.Nil(t, err)
	assert.NotNil(t, be)
	got, ok := cs.DecryptAES(secretKey, be)
	assert.Nil(t, ok)
	assert.Equal(t, expected, string(got))
}

func getTestServerCrypto() Crypto {
	crypt := new(crypto)
	crypt.privateKey =
		getTestConfigServerCrypto().
			Property().
			PrivateKey()
	return crypt
}

func init() {

	id := rand.Uint32()
	testPublicKeyFileName = fmt.Sprintf("%s/test_public_key_%018d.pem", os.TempDir(), id)
	testPrivateKeyFileName = fmt.Sprintf("%s/test_private_key_%018d.pem", os.TempDir(), id)
	privateKey, publicKey := GenerateRsaKeyPair()

	privateKeyFile, err := os.OpenFile(testPrivateKeyFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	util.IfErrorThenPanic(err)
	defer util.FileClose(privateKeyFile)
	_, err = privateKeyFile.Write([]byte(exportRsaPrivateKeyAsPemStr(privateKey)))
	util.IfErrorThenPanic(err)

	publicKeyStr, err := exportRsaPublicKeyAsPemStr(publicKey)
	util.IfErrorThenPanic(err)
	publicKeyFile, err := os.OpenFile(testPublicKeyFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	util.IfErrorThenPanic(err)
	defer util.FileClose(publicKeyFile)
	_, err = publicKeyFile.Write([]byte(publicKeyStr))
	util.IfErrorThenPanic(err)
}
