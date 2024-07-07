/*
 * This file was last modified at 2024-07-08 13:46 by Victor N. Skurikhin.
 * load_tls.go
 * $Id$
 */

package crypto

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

func LoadAgentTLSCredentials() (credentials.TransportCredentials, error) {
	// Загрузка сертификата центра сертификации, подписавшего сертификат сервера.
	pemServerCA, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	// Создание учётных данных для конфигурации TLS.
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}

func LoadServerTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	// Загрузка серверного сертификата и закрытого ключа.
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Создание учётных данных для конфигурации TLS.
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
