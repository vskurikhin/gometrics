/*
 * This file was last modified at 2024-06-15 16:00 by Victor N. Skurikhin.
 * simplersa.go
 * $Id$
 */

// Package simplersa Простая реализация шифрования RSA.
// Please use crypto/rsa in real code.
// Учебный пример, источник:
// https://kovardin.ru/articles/go/rsa/
package simplersa

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
)

// PublicKey структура для публичного ключа RSA.
// Размер для n не меньше 8192 бит в для агента см. TODO в agent.NewRequest.
type PublicKey struct {
	N *big.Int
	E *big.Int
}

// PrivateKey структура для приватного ключа.
// Согласно RFC 2313 N и D достаточны для расшифровки сообщений.
type PrivateKey struct {
	N *big.Int
	D *big.Int
}

// GenerateKeys генерирует пару ключей открытый/закрытый для шифрования/дешифрования RSA
// с заданной длиной. См. раздел 6 RFC 2313.
func GenerateKeys(bitlen int) (*PublicKey, *PrivateKey, error) {
	numRetries := 0

	for {
		numRetries++
		if numRetries == 10 {
			panic("retrying too many times, something is wrong")
		}

		// Результат p*q с b битами, поэтому мы генерируем p и q с b/2 битами
		// каждый. Если установлен старший бит p и q, результат будет содержать b битов.
		// В противном случае мы повторим попытку. rand.Prime должен возвращать простые числа
		// с установленным верхним битом, поэтому на практике повторных попыток не будет.
		p, err := rand.Prime(rand.Reader, bitlen/2)
		if err != nil {
			return nil, nil, err
		}
		q, err := rand.Prime(rand.Reader, bitlen/2)
		if err != nil {
			return nil, nil, err
		}

		// n это p*q
		n := new(big.Int).Set(p)
		n.Mul(n, q)

		if n.BitLen() != bitlen {
			continue
		}

		// theta(n) = (p-1)(q-1)
		p.Sub(p, big.NewInt(1))
		q.Sub(q, big.NewInt(1))
		totient := new(big.Int).Set(p)
		totient.Mul(totient, q)

		// e рекомендовано PKCS#1 (RFC 2313)
		e := big.NewInt(65537)

		// Вычислить обратную величину e по модулю умножения (p-1)*(q-1) = totient такую, что:
		// d*e = 1 (модульный коэффициент)
		// Если gcd(e, totient) = 1, то e гарантированно будет иметь уникальную инверсию,
		// но поскольку p-1 или q-1 теоретически могут иметь e в качестве множителя,
		// время от времени это может давать сбой (вероятно, это будет чрезвычайно редко).
		d := new(big.Int).ModInverse(e, totient)
		if d == nil {
			continue
		}

		pub := &PublicKey{N: n, E: e}
		priv := &PrivateKey{N: n, D: d}
		return pub, priv, nil
	}
}

// encrypt выполняет шифрование сообщения m с использованием открытого ключа
// и возвращает зашифрованное сообщение (big.Int).
// Кодирование сообщения как big.Int является обязанностью вызывающей стороны.
func encrypt(pub *PublicKey, m *big.Int) *big.Int {
	c := new(big.Int)
	c.Exp(m, pub.E, pub.N)
	return c
}

// decrypt расшифровывает зашифрованное сообщение (big.Int) с использованием
// закрытого ключа и возвращает расшифрованное сообщение.
func decrypt(priv *PrivateKey, c *big.Int) *big.Int {
	m := new(big.Int)
	m.Exp(c, priv.D, priv.N)
	return m
}

// EncryptRSA шифрует сообщение m с использованием открытого ключа pub
// и возвращает зашифрованные байты. Длина m должна быть <= size_in_bytes(pub.N) - 11,
// в противном случае возвращается ошибка.
// Формат блока шифрования основан на PKCS #1 v1.5 (RFC 2313).
func EncryptRSA(pub *PublicKey, m []byte) ([]byte, error) {
	// вычислить длину ключа в байтах, округляя вверх.
	keyLen := (pub.N.BitLen() + 7) / 8

	if len(m) > keyLen-11 {
		return nil, fmt.Errorf("len(m)=%v, too long", len(m))
	}

	// Согласно RFC 2313, для шифрования рекомендуется использовать тип блока 02:
	// EB = 00 || 02 || PS || 00 || D
	psLen := keyLen - len(m) - 3
	eb := make([]byte, keyLen)
	eb[0] = 0x00
	eb[1] = 0x02

	// Заполнить PS случайными ненулевыми байтами.
	for i := 2; i < 2+psLen; {
		_, err := rand.Read(eb[i : i+1])
		if err != nil {
			return nil, err
		}
		if eb[i] != 0x00 {
			i++
		}
	}
	eb[2+psLen] = 0x00

	// Скопировать сообщение m в оставшуюся часть блока шифрования.
	copy(eb[3+psLen:], m)

	// Блок шифрования; мы принимаем его как big.Int
	// и шифруем RSA с помощью открытого ключа.
	mnum := new(big.Int).SetBytes(eb)
	c := encrypt(pub, mnum)

	// В результате получается big.Int, который преобразовывается в байтовый фрагмент длиной keyLen.
	// Весьма, вероятно, что размер c в байтах равен keyLen, но в редких случаях нам может потребоваться
	// дополнить его слева нулями.
	// (это происходит только в том случае, если весь старший бит c равен нулям,
	// то есть он более чем в 256 раз меньше, чем модуль).
	padLen := keyLen - len(c.Bytes())
	for i := 0; i < padLen; i++ {
		eb[i] = 0x00
	}
	copy(eb[padLen:], c.Bytes())
	return eb, nil
}

// DecryptRSA расшифровывает сообщение c использованием закрытого ключа priv
// и возвращает расшифрованные байты на основе блока 02 из PKCS #1 v1.5 (RCS 2313).
// Он ожидает, что длина закрытого ключа в байтах по модулю будет равна len(eb).
// Важно: это простая реализация, не предназначенная для устойчивости к атакам по времени.
func DecryptRSA(priv *PrivateKey, c []byte) ([]byte, error) {
	keyLen := (priv.N.BitLen() + 7) / 8
	if len(c) != keyLen {
		return nil, fmt.Errorf("len(c)=%v, want keyLen=%v", len(c), keyLen)
	}

	// Преобразуйте c в bit.Int и расшифровывает его с помощью закрытого ключа.
	cnum := new(big.Int).SetBytes(c)
	mnum := decrypt(priv, cnum)

	// Записывает байты mnum в m, при необходимости заполняя их слева.
	m := make([]byte, keyLen)
	copy(m[keyLen-len(mnum.Bytes()):], mnum.Bytes())

	// Ожидает правильного начала блока 02.
	if m[0] != 0x00 {
		return nil, fmt.Errorf("m[0]=%v, want 0x00", m[0])
	}
	if m[1] != 0x02 {
		return nil, fmt.Errorf("m[1]=%v, want 0x02", m[1])
	}

	// Пропускает случайное заполнение, пока не будет достигнут байт 0x00.
	// +2 возвращает индекс к полному фрагменту.
	endPad := bytes.IndexByte(m[2:], 0x00) + 2
	if endPad < 2 {
		return nil, fmt.Errorf("end of padding not found")
	}

	return m[endPad+1:], nil
}
