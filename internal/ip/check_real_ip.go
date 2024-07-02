/*
 * This file was last modified at 2024-07-04 17:29 by Victor N. Skurikhin.
 * check_real_ip.go
 * $Id$
 */

// Package ip проверка вхождения в доверенную подсеть.
package ip

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/vskurikhin/gometrics/internal/env"
	"net"
	"net/http"
)

type HTTPError struct {
	Error string `json:"error"`
}

func (e *HTTPError) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

var ErrUserUnauthorized = fmt.Errorf("ip не авторизован")

func XRealIPChecker(next http.Handler) http.Handler {
	return xRealIPChecker(next, ErrUserUnauthorized)
}

func xRealIPChecker(next http.Handler, _ error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := env.GetServerConfig()
		trustedSubnet := cfg.TrustedSubnet()
		if trustedSubnet != "" {

			xRealIP := r.Header.Get("X-Real-IP")
			_, ipNet, err := net.ParseCIDR(trustedSubnet)

			if err != nil {
				http.Error(w, "", http.StatusForbidden)
				//goland:noinspection GoUnhandledErrorResult
				_ = render.Render(w, r, &HTTPError{Error: err.Error()})
			} else if xRealIP == "" {
				http.Error(w, "", http.StatusForbidden)
				//goland:noinspection GoUnhandledErrorResult
				_ = render.Render(w, r, &HTTPError{Error: "empty X-Real-IP"})
			} else if !ipNet.Contains(net.ParseIP(xRealIP)) {
				http.Error(w, "", http.StatusForbidden)
				//goland:noinspection GoUnhandledErrorResult
				_ = render.Render(w, r, &HTTPError{Error: "forbidden for X-Real-IP"})
			} else {
				next.ServeHTTP(w, r)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
