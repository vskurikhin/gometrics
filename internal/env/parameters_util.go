/*
 * This file was last modified at 2024-07-04 17:29 by Victor N. Skurikhin.
 * parameters_util.go
 * $Id$
 */

package env

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net"
	"os"
)

// getOutboundIP - Get preferred outbound ip of this machine.
func getOutboundIP() net.IP {

	if env == nil {
		return nil
	}
	conn, err := net.Dial("udp", env.DNS)

	if err != nil {
		return nil
	}
	defer func() { _ = conn.Close() }()

	if localAddr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		return localAddr.IP
	}
	return nil
}

func loadPrivateKey() *rsa.PrivateKey {
	if len(cfg.CryptoKey()) > 1 {
		file, err := os.Open(cfg.CryptoKey())
		if err != nil {
			return nil
		}
		//nolint:multichecker,errcheck
		defer func() { _ = file.Close() }()
		buf, err := io.ReadAll(file)
		if err != nil {
			return nil
		}
		if block := readPEMString(string(buf)); block != nil {
			privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil
			}
			return privateKey
		}
	}
	return nil
}

func loadPublicKey() *rsa.PublicKey {
	if len(cfg.CryptoKey()) > 1 {
		file, err := os.Open(cfg.CryptoKey())
		if err != nil {
			return nil
		}
		//nolint:multichecker,errcheck
		defer func() { _ = file.Close() }()
		buf, err := io.ReadAll(file)
		if err != nil {
			return nil
		}
		if block := readPEMString(string(buf)); block != nil {
			publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return nil
			}
			return publicKey
		}
	}
	return nil
}

func readPEMString(p string) *pem.Block {
	result, _ := pem.Decode([]byte(p))
	return result
}
