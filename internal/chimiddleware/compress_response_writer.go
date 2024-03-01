/*
 * This file was last modified at 2024-03-01 21:38 by Victor N. Skurikhin.
 * compress_response_writer.go
 * $Id$
 */

package chimiddleware

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

// An EncoderFunc is a function that wraps the provided io.Writer with a
// streaming compression algorithm and returns it.
//
// In case of failure, the function should return nil.
type EncoderFunc func(w io.Writer, level int) io.Writer

type compressResponseWriter struct {
	http.ResponseWriter

	// The streaming encoder writer to be used if there is one. Otherwise,
	// this is just the normal writer.
	w                io.Writer
	encoding         string
	contentTypes     map[string]struct{}
	contentWildcards map[string]struct{}
	wroteHeader      bool
	compressable     bool
}

func (cw *compressResponseWriter) isCompressable() bool {
	// Parse the first part of the Content-Type response header.
	contentType := cw.Header().Get("Content-Type")
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = contentType[0:idx]
	}

	// Is the content type compressable?
	if _, ok := cw.contentTypes[contentType]; ok {
		return true
	}
	if idx := strings.Index(contentType, "/"); idx > 0 {
		contentType = contentType[0:idx]
		_, ok := cw.contentWildcards[contentType]
		return ok
	}
	return false
}

func (cw *compressResponseWriter) WriteHeader(code int) {
	if cw.wroteHeader {
		cw.ResponseWriter.WriteHeader(code) // Allow multiple calls to propagate.
		return
	}
	cw.wroteHeader = true
	defer cw.ResponseWriter.WriteHeader(code)

	// Already compressed data?
	if cw.Header().Get("Content-Encoding") != "" {
		return
	}

	if !cw.isCompressable() {
		cw.compressable = false
		return
	}

	if cw.encoding != "" {
		cw.compressable = true
		cw.Header().Set("Content-Encoding", cw.encoding)
		cw.Header().Set("Vary", "Accept-Encoding")

		// The content-length after compression is unknown
		cw.Header().Del("Content-Length")
	}
}

func (cw *compressResponseWriter) Write(p []byte) (int, error) {
	if !cw.wroteHeader {
		cw.WriteHeader(http.StatusOK)
	}

	return cw.writer().Write(p)
}

func (cw *compressResponseWriter) writer() io.Writer {
	if cw.compressable {
		return cw.w
	} else {
		return cw.ResponseWriter
	}
}

type compressFlusher interface {
	Flush() error
}

//goland:noinspection ALL
func (cw *compressResponseWriter) Flush() {
	if f, ok := cw.writer().(http.Flusher); ok {
		f.Flush()
	}
	// If the underlying writer has a compression flush signature,
	// call this Flush() method instead
	if f, ok := cw.writer().(compressFlusher); ok {
		f.Flush()

		// Also flush the underlying response writer
		if f, ok := cw.ResponseWriter.(http.Flusher); ok {
			f.Flush()
		}
	}
}

func (cw *compressResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := cw.writer().(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("chi/middleware: http.Hijacker is unavailable on the writer")
}

func (cw *compressResponseWriter) Push(target string, opts *http.PushOptions) error {
	if ps, ok := cw.writer().(http.Pusher); ok {
		return ps.Push(target, opts)
	}
	return errors.New("chi/middleware: http.Pusher is unavailable on the writer")
}

func (cw *compressResponseWriter) Close() error {
	if c, ok := cw.writer().(io.WriteCloser); ok {
		return c.Close()
	}
	return errors.New("chi/middleware: io.WriteCloser is unavailable on the writer")
}

func encoderGzip(w io.Writer, level int) io.Writer {
	gw, err := gzip.NewWriterLevel(w, level)
	if err != nil {
		return nil
	}
	return gw
}
