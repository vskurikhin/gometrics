/*
 * This file was last modified at 2024-07-04 17:29 by Victor N. Skurikhin.
 * parameters_util_test.go
 * $Id$
 */

package env

import (
	"crypto/rsa"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vskurikhin/gometrics/internal/util"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	publicKeyPEM = `-----BEGIN RSA PUBLIC KEY-----
MCgCIQDjGY/39sSUndnI5PnXKvruj+jOMOnlKuNun1sx9/npzQIDAQAB
-----END RSA PUBLIC KEY-----`
	privateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIGrAgEAAiEA4xmP9/bElJ3ZyOT51yr67o/ozjDp5Srjbp9bMff56c0CAwEAAQIg
aESyTz0joMCg35YSB/KZ5tTgiY42Ber6gkfop/7Kf9kCEQD7DQgEp3PiQIg3nfua
L5v7AhEA55Oogh+6xANRG0wGeNue1wIRAMDLRCIW2rag6jsT9vl0oGsCEQCMrUJ8
adIHKQyoTHLSEHhZAhA2af+67sFM8Vk7I2W0SPzF
-----END RSA PRIVATE KEY-----`
)

var (
	testPublicKeyFileName  string
	testPrivateKeyFileName string
)

func getTestConfigAgent() Config {
	return GetTestConfig(
		WithCryptoKey(testPublicKeyFileName),
	)
}

func getTestConfigServer() Config {
	return GetTestConfig(
		WithCryptoKey(testPrivateKeyFileName),
	)
}

func TestGetOutboundIP(t *testing.T) {
	gotNil := getOutboundIP()
	assert.Nil(t, gotNil)
	env = nil
	gotNil = getOutboundIP()
	assert.Nil(t, gotNil)
	onceEnv = new(sync.Once)
	getEnvironments()
	got := getOutboundIP()
	assert.NotNil(t, got)
	onceEnv = new(sync.Once)
	env = nil
}

func TestLoadPrivateKey(t *testing.T) {
	getTestConfigServer()
	n, ok1 := big.NewInt(0).SetString("e3198ff7f6c4949dd9c8e4f9d72afaee8fe8ce30e9e52ae36e9f5b31f7f9e9cd", 16)
	assert.True(t, ok1)
	publicKey := &rsa.PublicKey{E: 65537, N: n}
	privateKey := loadPrivateKey()
	expectedPrivateKey := privateKey
	expectedPrivateKey.N = publicKey.N
	expectedPrivateKey.E = publicKey.E
	d, ok2 := big.NewInt(0).SetString("6844b24f3d23a0c0a0df961207f299e6d4e0898e3605eafa8247e8a7feca7fd9", 16)
	assert.True(t, ok2)
	expectedPrivateKey.D = d
	p1, ok3 := big.NewInt(0).SetString("fb0d0804a773e24088379dfb9a2f9bfb", 16)
	assert.True(t, ok3)
	p2, ok4 := big.NewInt(0).SetString("e793a8821fbac403511b4c0678db9ed7", 16)
	assert.True(t, ok4)
	expectedPrivateKey.Primes = []*big.Int{p1, p2}
	dp, ok5 := big.NewInt(0).SetString("c0cb442216dab6a0ea3b13f6f974a06b", 16)
	assert.True(t, ok5)
	expectedPrivateKey.Precomputed.Dp = dp
	dq, ok6 := big.NewInt(0).SetString("8cad427c69d207290ca84c72d2107859", 16)
	assert.True(t, ok6)
	expectedPrivateKey.Precomputed.Dq = dq
	Qinv, ok7 := big.NewInt(0).SetString("3669ffbaeec14cf1593b2365b448fcc5", 16)
	assert.True(t, ok7)
	expectedPrivateKey.Precomputed.Qinv = Qinv
	assert.Equal(t, expectedPrivateKey, privateKey)
}

func TestLoadPublicKey(t *testing.T) {
	parameters = nil
	onceParameters = new(sync.Once)
	getTestConfigAgent()
	time.Sleep(501 * time.Millisecond)
	publicKey := loadPublicKey()
	n, ok1 := big.NewInt(0).SetString("e3198ff7f6c4949dd9c8e4f9d72afaee8fe8ce30e9e52ae36e9f5b31f7f9e9cd", 16)
	assert.True(t, ok1)
	expectedPublicKey := &rsa.PublicKey{E: 65537, N: n}
	assert.Equal(t, expectedPublicKey, publicKey)
}

func init() {

	id := rand.Uint32()
	testPublicKeyFileName = fmt.Sprintf("%s/test_public_key_%018d.pem", os.TempDir(), id)
	testPrivateKeyFileName = fmt.Sprintf("%s/test_private_key_%018d.pem", os.TempDir(), id)
	privateKeyFile, err := os.OpenFile(testPrivateKeyFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	util.IfErrorThenPanic(err)
	defer func() { _ = privateKeyFile.Close() }()
	_, err = privateKeyFile.Write([]byte(privateKeyPEM))
	util.IfErrorThenPanic(err)

	publicKeyFile, err := os.OpenFile(testPublicKeyFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	util.IfErrorThenPanic(err)
	defer func() { _ = publicKeyFile.Close() }()
	_, err = publicKeyFile.Write([]byte(publicKeyPEM))
	util.IfErrorThenPanic(err)
}
