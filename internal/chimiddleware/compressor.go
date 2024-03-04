/*
 * This file was last modified at 2024-03-02 12:35 by Victor N. Skurikhin.
 * compressor.go
 * $Id$
 */

package chimiddleware

import (
	"io"
	"net/http"
	"strings"
	"sync"
)

// NewCompressor creates a new Compressor that will handle encoding responses.
//
// The level should be one of the ones defined in the flate package.
// The types are the content types that are allowed to be compressed.
func NewCompressor(level int, types ...string) *Compressor {
	// If types are provided, set those as the allowed types. If none are
	// provided, use the default list.
	allowedTypes := make(map[string]struct{})
	allowedWildcards := make(map[string]struct{})
	for _, t := range types {
		allowedTypes[t] = struct{}{}
	}

	c := &Compressor{
		level:            level,
		encoders:         make(map[string]EncoderFunc),
		pooledEncoders:   make(map[string]*sync.Pool),
		allowedTypes:     allowedTypes,
		allowedWildcards: allowedWildcards,
	}

	// TODO: Exception for old MSIE browsers that can't handle non-HTML?
	// https://zoompf.com/blog/2012/02/lose-the-wait-http-compression
	c.SetEncoder("gzip", encoderGzip)

	return c
}

// Compressor represents a set of encoding configurations.
type Compressor struct {
	level int // The compression level.
	// The mapping of encoder names to encoder functions.
	encoders map[string]EncoderFunc
	// The mapping of pooled encoders to pools.
	pooledEncoders map[string]*sync.Pool
	// The set of content types allowed to be compressed.
	allowedTypes     map[string]struct{}
	allowedWildcards map[string]struct{}
	// The list of encoders in order of decreasing precedence.
	encodingPrecedence []string
}

// SetEncoder can be used to set the implementation of a compression algorithm.
//
// The encoding should be a standardised identifier. See:
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Encoding
//
// For example, add the Brotli algortithm:
//
//	import brotli_enc "gopkg.in/kothar/brotli-go.v0/enc"
//
//	compressor := middleware.NewCompressor(5, "text/html")
//	compressor.SetEncoder("br", func(w http.ResponseWriter, level int) io.Writer {
//	  params := brotli_enc.NewBrotliParams()
//	  params.SetQuality(level)
//	  return brotli_enc.NewBrotliWriter(params, w)
//	})
func (c *Compressor) SetEncoder(encoding string, fn EncoderFunc) {
	encoding = strings.ToLower(encoding)
	if encoding == "" {
		panic("the encoding can not be empty")
	}
	if fn == nil {
		panic("attempted to set a nil encoder function")
	}

	// If we are adding a new encoder that is already registered, we have to
	// clear that one out first.
	delete(c.pooledEncoders, encoding)
	delete(c.encoders, encoding)

	// If the encoder supports Resetting (IoReseterWriter), then it can be pooled.
	encoder := fn(io.Discard, c.level)
	if _, ok := encoder.(ioResetterWriter); ok {
		pool := &sync.Pool{
			New: func() interface{} {
				return fn(io.Discard, c.level)
			},
		}
		c.pooledEncoders[encoding] = pool
	}
	// If the encoder is not in the pooledEncoders, add it to the normal encoders.
	if _, ok := c.pooledEncoders[encoding]; !ok {
		c.encoders[encoding] = fn
	}

	for i, v := range c.encodingPrecedence {
		if v == encoding {
			c.encodingPrecedence = append(c.encodingPrecedence[:i], c.encodingPrecedence[i+1:]...)
		}
	}

	c.encodingPrecedence = append([]string{encoding}, c.encodingPrecedence...)
}

// Handler returns a new middleware that will compress the response based on the
// current Compressor.
//
//goland:noinspection GoUnhandledErrorResult
func (c *Compressor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder, encoding, cleanup := c.selectEncoder(r.Header, w)

		cw := &compressResponseWriter{
			ResponseWriter:   w,
			w:                w,
			contentTypes:     c.allowedTypes,
			contentWildcards: c.allowedWildcards,
			encoding:         encoding,
			compressable:     false, // determined in post-handler
		}
		if encoder != nil {
			cw.w = encoder
		}
		// Re-add the encoder to the pool if applicable.
		defer cleanup()
		defer cw.Close()

		next.ServeHTTP(cw, r)
	})
}

// selectEncoder returns the encoder, the name of the encoder, and a closer function.
func (c *Compressor) selectEncoder(h http.Header, w io.Writer) (io.Writer, string, func()) {
	header := h.Get("Accept-Encoding")

	// Parse the names of all accepted algorithms from the header.
	accepted := strings.Split(strings.ToLower(header), ",")

	// Find supported encoder by accepted list by precedence
	for _, name := range c.encodingPrecedence {
		if matchAcceptEncoding(accepted, name) {
			if pool, ok := c.pooledEncoders[name]; ok {
				encoder := pool.Get().(ioResetterWriter)
				cleanup := func() {
					pool.Put(encoder)
				}
				encoder.Reset(w)
				return encoder, name, cleanup

			}
			if fn, ok := c.encoders[name]; ok {
				return fn(w, c.level), name, func() {}
			}
		}

	}

	// No encoder found to match the accepted encoding
	return nil, "", func() {}
}

func matchAcceptEncoding(accepted []string, encoding string) bool {
	for _, v := range accepted {
		if strings.Contains(v, encoding) {
			return true
		}
	}
	return false
}

// Interface for types that allow resetting io.Writers.
type ioResetterWriter interface {
	io.Writer
	Reset(w io.Writer)
}
