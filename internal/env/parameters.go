/*
 * This file was last modified at 2024-07-04 17:29 by Victor N. Skurikhin.
 * parameters.go
 * $Id$
 */

package env

import (
	"crypto/rsa"
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	OutboundIP = "OUTBOUND_IP"
	PrivateKey = "PRIVATE_KEY"
	PublicKey  = "PUBLIC_KEY"
)

type Parameters interface {
	fmt.Stringer
	OutboundIP() net.IP
	PrivateKey() *rsa.PrivateKey
	PublicKey() *rsa.PublicKey
}

type mapParameters struct {
	mp sync.Map
}

var _ Parameters = (*mapParameters)(nil)
var onceParameters = new(sync.Once)
var parameters *mapParameters

// GetParameters — параметры преобразованные из конфигурации и окружения.
func GetParameters() Parameters {

	onceParameters.Do(func() {
		parameters = new(mapParameters)
		go func() {
			for {
				time.Sleep(500 * time.Millisecond)
				if ip := getOutboundIP(); ip != nil {
					parameters.mp.Store(OutboundIP, ip)
				}
			}
		}()
	})
	return parameters
}

func setParameters() {

	GetParameters()
	parameters.mp.Store(OutboundIP, net.IPv4(127, 0, 0, 1))

	if privateKey := loadPrivateKey(); privateKey != nil {
		parameters.mp.Store(PrivateKey, privateKey)
	}
	if publicKey := loadPublicKey(); publicKey != nil {
		parameters.mp.Store(PublicKey, publicKey)
	}
}

func (p *mapParameters) OutboundIP() net.IP {

	if ip, ok := p.mp.Load(OutboundIP); ok {
		return ip.(net.IP)
	}
	return nil
}

func (p *mapParameters) PrivateKey() *rsa.PrivateKey {

	if ip, ok := p.mp.Load(PrivateKey); ok {
		return ip.(*rsa.PrivateKey)
	}
	return nil
}

func (p *mapParameters) PublicKey() *rsa.PublicKey {

	if ip, ok := p.mp.Load(PublicKey); ok {
		return ip.(*rsa.PublicKey)
	}
	return nil
}

func (p *mapParameters) String() string {
	format := `
	OutboundIP : %v
	PrivateKey : %v
	PublicKey  : %v
`
	return fmt.Sprintf(format,
		p.OutboundIP(),
		p.PrivateKey(),
		p.PublicKey(),
	)
}
