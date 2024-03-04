/*
 * This file was last modified at 2024-03-01 21:41 by Victor N. Skurikhin.
 * compress.go
 * $Id$
 */

package compress

import (
	"github.com/vskurikhin/gometrics/internal/chimiddleware"
	"net/http"
)

var defaultLevel = 5

var defaultCompressibleContentTypes = []string{
	"text/html",
	"application/json",
}

// Compress returns a logger handler.
func Compress(next http.Handler) http.Handler {
	compressor := chimiddleware.NewCompressor(defaultLevel, defaultCompressibleContentTypes...)
	return compressor.Handler(next)
}
