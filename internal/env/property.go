/*
 * This file was last modified at 2024-07-08 14:02 by Victor N. Skurikhin.
 * property.go
 * $Id$
 */

package env

import (
	"crypto/rsa"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vskurikhin/gometrics/internal/storage"
)

const (
	DBPool     = "DB_POOL"
	OutboundIP = "OUTBOUND_IP"
	PrivateKey = "PRIVATE_KEY"
	PublicKey  = "PUBLIC_KEY"
	Storage    = "STORAGE"
)

type Property interface {
	fmt.Stringer
	DBPool() *pgxpool.Pool
	OutboundIP() net.IP
	PrivateKey() *rsa.PrivateKey
	PublicKey() *rsa.PublicKey
	Storage() storage.Storage
}

type mapProperty struct {
	mp sync.Map
}

var _ Property = (*mapProperty)(nil)
var onceProperty = new(sync.Once)

// GetProperty — свойства преобразованные из конфигурации и окружения.
func GetProperty() Property {

	var property *mapProperty
	var pool *pgxpool.Pool
	if cfg.IsDBSetup() {
		pool = dbConnect(cfg.DataBaseDSN())
	}
	property = getProperty(
		withDBPool(pool),
		WithPrivateKey(LoadPrivateKey()),
		WithPublicKey(LoadPublicKey()),
		WithOutboundIP(net.IPv4(127, 0, 0, 1)),
		WithStorage(getStorage(pool)),
	)
	onceProperty.Do(func() {
		go func() {
			for {
				time.Sleep(500 * time.Millisecond)
				if ip := getOutboundIP(); ip != nil {
					property.mp.Store(OutboundIP, ip)
				}
			}
		}()
	})
	cfg.property = property
	return cfg.property
}

// GetTestProperty — для создания тестовой конфигурации.
func GetTestProperty(opts ...func(*mapProperty)) Property {
	return getProperty(opts...)
}

// withDBPool — пул соединений с БД.
func withDBPool(pool *pgxpool.Pool) func(*mapProperty) {
	return func(p *mapProperty) {
		if pool != nil {
			p.mp.Store(DBPool, pool)
		}
	}
}

func (p *mapProperty) DBPool() *pgxpool.Pool {

	if dbPool, ok := p.mp.Load(DBPool); ok {
		return dbPool.(*pgxpool.Pool)
	}
	return nil
}

// WithOutboundIP — исходящий IP адрес.
func WithOutboundIP(ip net.IP) func(*mapProperty) {
	return func(p *mapProperty) {
		if ip != nil {
			p.mp.Store(OutboundIP, ip)
		}
	}
}

func (p *mapProperty) OutboundIP() net.IP {

	if ip, ok := p.mp.Load(OutboundIP); ok {
		return ip.(net.IP)
	}
	return nil
}

// WithPrivateKey — секретный ключ.
func WithPrivateKey(privateKey *rsa.PrivateKey) func(*mapProperty) {
	return func(p *mapProperty) {
		if privateKey != nil {
			p.mp.Store(PrivateKey, privateKey)
		}
	}
}

func (p *mapProperty) PrivateKey() *rsa.PrivateKey {

	if ip, ok := p.mp.Load(PrivateKey); ok {
		return ip.(*rsa.PrivateKey)
	}
	return nil
}

// WithPublicKey — публичный ключ.
func WithPublicKey(publicKey *rsa.PublicKey) func(*mapProperty) {
	return func(p *mapProperty) {
		if publicKey != nil {
			p.mp.Store(PublicKey, publicKey)
		}
	}
}

func (p *mapProperty) PublicKey() *rsa.PublicKey {

	if ip, ok := p.mp.Load(PublicKey); ok {
		return ip.(*rsa.PublicKey)
	}
	return nil
}

// WithStorage — подсистема хранения.
func WithStorage(store storage.Storage) func(*mapProperty) {
	return func(p *mapProperty) {
		if store != nil {
			p.mp.Store(Storage, store)
		}
	}
}

func (p *mapProperty) Storage() storage.Storage {

	if store, ok := p.mp.Load(Storage); ok {
		return store.(storage.Storage)
	}
	return nil
}

func (p *mapProperty) String() string {
	format := `
	OutboundIP : %v
	PrivateKey : %v
	PublicKey  : %v
    Storage    : %v
`
	return fmt.Sprintf(format,
		p.OutboundIP(),
		p.PrivateKey(),
		p.PublicKey(),
		p.Storage(),
	)
}

func getProperty(opts ...func(*mapProperty)) *mapProperty {

	var property = new(mapProperty)

	// вызываем все указанные функции для установки параметров
	for _, opt := range opts {
		opt(property)
	}

	return property
}
