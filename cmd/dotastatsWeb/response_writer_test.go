package main_test

import (
	"bufio"
	main "dotastats/cmd/dotastatsWeb"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type hijackableResponse struct {
	Hijacked bool
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{}
}

func (h *hijackableResponse) Header() http.Header           { return nil }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func TestResponseWriterWritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := main.NewResponseWriter(rec)

	_, err := rw.Write([]byte("Hello world"))
	ok(t, err)
	equals(t, rec.Code, rw.Status())
	equals(t, rec.Body.String(), "Hello world")
	equals(t, rw.Status(), http.StatusOK)
	equals(t, rw.Size(), 11)
	equals(t, rw.Written(), true)
}

func TestResponseWriterWritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := main.NewResponseWriter(rec)

	_, err := rw.Write([]byte("Hello world"))
	ok(t, err)
	_, err = rw.Write([]byte("foo bar bat baz"))
	ok(t, err)

	equals(t, rec.Code, rw.Status())
	equals(t, rec.Body.String(), "Hello worldfoo bar bat baz")
	equals(t, rw.Status(), http.StatusOK)
	equals(t, rw.Size(), 26)
}

func TestResponseWriterWritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := main.NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)

	equals(t, rec.Code, rw.Status())
	equals(t, rec.Body.String(), "")
	equals(t, rw.Status(), http.StatusNotFound)
	equals(t, rw.Size(), 0)
}

func TestResponseWriterBefore(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := main.NewResponseWriter(rec)
	result := ""

	rw.Before(func(main.ResponseWriter) {
		result += "foo"
	})
	rw.Before(func(main.ResponseWriter) {
		result += "bar"
	})

	rw.WriteHeader(http.StatusNotFound)

	equals(t, rec.Code, rw.Status())
	equals(t, rec.Body.String(), "")
	equals(t, rw.Status(), http.StatusNotFound)
	equals(t, rw.Size(), 0)
	equals(t, result, "barfoo")
}

func TestResponseWriterHijack(t *testing.T) {
	hijackable := newHijackableResponse()
	rw := main.NewResponseWriter(hijackable)
	hijacker, ok := rw.(http.Hijacker)
	equals(t, ok, true)
	_, _, err := hijacker.Hijack()
	if err != nil {
		t.Error(err)
	}
	equals(t, hijackable.Hijacked, true)
}

func TestResponseWriteHijackNotOK(t *testing.T) {
	hijackable := new(http.ResponseWriter)
	rw := main.NewResponseWriter(*hijackable)
	hijacker, ok := rw.(http.Hijacker)
	equals(t, ok, true)
	_, _, err := hijacker.Hijack()

	assert(t, err != nil, "err is supposed to be returned")
}

func TestResponseWriterCloseNotify(t *testing.T) {
	rec := newCloseNotifyingRecorder()
	rw := main.NewResponseWriter(rec)
	closed := false
	notifier := rw.(http.CloseNotifier).CloseNotify()
	rec.close()
	select {
	case <-notifier:
		closed = true
	case <-time.After(time.Second):
	}
	equals(t, closed, true)
}

func TestResponseWriterFlusher(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := main.NewResponseWriter(rec)

	_, ok := rw.(http.Flusher)
	equals(t, ok, true)
}
