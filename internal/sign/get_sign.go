/*
 * This file was last modified at 2024-04-05 08:40 by Victor N. Skurikhin.
 * get_sign.go
 * $Id$
 */

package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/vskurikhin/gometrics/internal/env"
	"io"
	"net/http"
)

func GetSignFromRequest(request *http.Request) (string, error) {

	buf, err := io.ReadAll(request.Body)

	if err != nil {
		return "", err
	}
	r := io.NopCloser(bytes.NewBuffer(buf))
	request.Body = io.NopCloser(bytes.NewBuffer(buf))

	body, _ := io.ReadAll(r)
	h := hmac.New(sha256.New, []byte(*env.Server.Key()))
	h.Write(body)
	signSum := h.Sum(nil)

	return fmt.Sprintf("%x", signSum), nil
}
