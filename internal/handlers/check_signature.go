/*
 * This file was last modified at 2024-04-05 10:14 by Victor N. Skurikhin.
 * check_signature.go
 * $Id$
 */

package handlers

import (
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"github.com/vskurikhin/gometrics/internal/logger"
	"github.com/vskurikhin/gometrics/internal/sign"
	"go.uber.org/zap"
	"net/http"
)

func checkSignature(request *http.Request) (int, error) {

	if *env.Server.Key() != "" {
		dst, err := sign.GetSignatureFromRequest(request)
		if err != nil {
			return http.StatusBadRequest, err
		}
		hashSHA256 := request.Header.Get("HashSHA256")

		if hashSHA256 != dst {
			logger.Log.Debug("signature is bad", zap.String("HashSHA256", fmt.Sprintf("%x", hashSHA256)))
			return http.StatusBadRequest, fmt.Errorf("signature is bad: %s", hashSHA256)
		}
	}
	return http.StatusOK, nil
}
