/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * basic_writer.go
 * $Id$
 */

package chimiddleware

import (
	"io"
	"net/http"
)

// basicWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
	tee         io.Writer
}

func (b *basicWriter) WriteHeader(code int) {
	if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *basicWriter) Write(buf []byte) (int, error) {
	b.maybeWriteHeader()
	n, err := b.ResponseWriter.Write(buf)
	if b.tee != nil {
		_, err2 := b.tee.Write(buf[:n])
		// Prefer errors generated by the proxied writer.
		if err == nil {
			err = err2
		}
	}
	b.bytes += n
	return n, err
}

func (b *basicWriter) maybeWriteHeader() {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}
}

func (b *basicWriter) Status() int {
	return b.code
}

func (b *basicWriter) BytesWritten() int {
	return b.bytes
}

func (b *basicWriter) Tee(w io.Writer) {
	b.tee = w
}

func (b *basicWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}