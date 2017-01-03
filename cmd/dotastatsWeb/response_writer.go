package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWrite that provides extra infomation about
// The response. It is recommanded that middleware handlers user this construct to wrap a responsewriter
// If the functionality calls for it.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 0 if the response has not bean written.
	Status() int
	// Written returns whether or not the ResponseWriter has bean written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
	// Before allows for a function to be called before the ResponseWritter has been written to. This is
	// usefull for setting headers or any other operations that must happen before a response has been written
	Before(func(ResponseWriter))
}

type beforeFunc func(ResponseWriter)

// NewResponseWritter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0, 0, nil}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []beforeFunc
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.status = s
	rw.callBefore()
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Before(before func(ResponseWriter)) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("The ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
