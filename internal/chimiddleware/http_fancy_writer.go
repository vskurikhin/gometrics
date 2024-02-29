/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * http_fancy_writer.go
 * $Id$
 */

package chimiddleware

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

// http2FancyWriter is a HTTP2 writer that additionally satisfies
// http.Flusher, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to
// make the proxied object support the full method set of the proxied object.
type http2FancyWriter struct {
	basicWriter
}

func (f *http2FancyWriter) Flush() {
	f.wroteHeader = true
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

var _ http.Flusher = &http2FancyWriter{}

// httpFancyWriter is a HTTP writer that additionally satisfies
// http.Flusher, http.Hijacker, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to
// make the proxied object support the full method set of the proxied object.
type httpFancyWriter struct {
	basicWriter
}

func (f *httpFancyWriter) Flush() {
	f.wroteHeader = true
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (f *httpFancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (f *http2FancyWriter) Push(target string, opts *http.PushOptions) error {
	return f.basicWriter.ResponseWriter.(http.Pusher).Push(target, opts)
}

func (f *httpFancyWriter) ReadFrom(r io.Reader) (int64, error) {
	if f.basicWriter.tee != nil {
		n, err := io.Copy(&f.basicWriter, r)
		f.basicWriter.bytes += int(n)
		return n, err
	}
	rf := f.basicWriter.ResponseWriter.(io.ReaderFrom)
	f.basicWriter.maybeWriteHeader()
	n, err := rf.ReadFrom(r)
	f.basicWriter.bytes += int(n)
	return n, err
}

var _ http.Flusher = &httpFancyWriter{}
var _ http.Hijacker = &httpFancyWriter{}
var _ http.Pusher = &http2FancyWriter{}
var _ io.ReaderFrom = &httpFancyWriter{}
