/*
 * This file was last modified at 2024-02-04 14:23 by Victor N. Skurikhin.
 * format_request.go
 * $Id$
 */

package util

import (
	"fmt"
	"net/http"
	"strings"
)

// FormatRequest генерирует ascii-представление запроса
func FormatRequest(r *http.Request) string {

	// Создать возвращаемую строку
	var request []string // Добавить строку запроса

	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)                             // Добавить хост
	request = append(request, fmt.Sprintf("Host: %v", r.Host)) // Перебирать заголовки

	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// Если это POST, добавить данные
	if r.Method == http.MethodPost {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	return strings.Join(request, "\n")
}
