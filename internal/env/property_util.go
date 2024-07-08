/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * property_util.go
 * $Id$
 */

package env

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/storage"
	"github.com/vskurikhin/gometrics/internal/util"
)

func dbConnect(dsn string) *pgxpool.Pool {

	config, err := pgxpool.ParseConfig(dsn)
	util.IfErrorThenPanic(err)
	logger.Log.Debug("dbConnect config parsed")

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		logger.Log.Debug("Acquire connect ping...")
		if err = conn.Ping(ctx); err != nil {
			panic(err)
		}
		logger.Log.Debug("Acquire connect Ok")
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.TODO(), config)
	util.IfErrorThenPanic(err)
	logger.Log.Debug("NewWithConfig pool created")
	_, err = pool.Acquire(context.TODO())
	util.IfErrorThenPanic(err)
	logger.Log.Debug("Acquire pool Ok")

	return pool
}

func getStorage(pool *pgxpool.Pool) storage.Storage {

	mem := new(storage.MemStorage)
	mem.Metrics = make(map[string]*string)

	if cfg.IsDBSetup() {
		return storage.New(mem, pool)
	} else {
		return mem
	}
}

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

func LoadPrivateKey() *rsa.PrivateKey {
	if len(cfg.CryptoKey()) > 1 {
		file, err := os.Open(cfg.CryptoKey())
		if err != nil {
			return nil
		}
		//nolint:multichecker,errcheck
		defer util.FileClose(file)
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

func LoadPublicKey() *rsa.PublicKey {
	if len(cfg.CryptoKey()) > 1 {
		file, err := os.Open(cfg.CryptoKey())
		if err != nil {
			return nil
		}
		//nolint:multichecker,errcheck
		defer util.FileClose(file)
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
