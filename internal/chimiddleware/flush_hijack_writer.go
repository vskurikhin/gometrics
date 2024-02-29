/*
 * This file was last modified at 2024-02-29 12:49 by Victor N. Skurikhin.
 * flush_hijack_writer.go
 * $Id$
 */

package chimiddleware

import (
	"bufio"
	"net"
	"net/http"
)

// hijackWriter ...
type hijackWriter struct {
	basicWriter
}

func (f *hijackWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

var _ http.Hijacker = &hijackWriter{}

// flushHijackWriter ...
type flushHijackWriter struct {
	basicWriter
}

func (f *flushHijackWriter) Flush() {
	f.wroteHeader = true
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

func (f *flushHijackWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := f.basicWriter.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

var _ http.Flusher = &flushHijackWriter{}
var _ http.Hijacker = &flushHijackWriter{}

// flushWriter ...
type flushWriter struct {
	basicWriter
}

func (f *flushWriter) Flush() {
	f.wroteHeader = true
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}

var _ http.Flusher = &flushWriter{}
