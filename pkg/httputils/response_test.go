package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCodeLoggingResponseWriter(t *testing.T) {
	t.Parallel()

	rw := httptest.NewRecorder()
	scrw := NewStatusCodeResponseWriter(rw)

	assert.Equal(t, rw, scrw.wrapped)
	assert.Implements(t, (*http.ResponseWriter)(nil), scrw)
}

func TestStatusCodeResponseWriter_Header(t *testing.T) {
	t.Parallel()

	const (
		key    = "Some-Header"
		value1 = "value1"
		value2 = "value2"
	)

	headers := http.Header{
		key: {value1, value2},
	}

	rw := httptest.NewRecorder()
	rw.Header().Add(key, value1)
	rw.Header().Add(key, value2)

	scrw := StatusCodeResponseWriter{wrapped: rw}

	assert.Equal(t, headers, scrw.Header())
}

func TestStatusCodeResponseWriter_StatusCode(t *testing.T) {
	t.Parallel()

	const status = http.StatusTeapot

	scrw := StatusCodeResponseWriter{statusCode: status}

	assert.Equal(t, status, scrw.StatusCode())
}

func TestStatusCodeResponseWriter_Write(t *testing.T) {
	t.Parallel()

	content := []byte("some-content")

	rw := httptest.NewRecorder()
	scrw := &StatusCodeResponseWriter{wrapped: rw}

	n, err := scrw.Write(content)

	require.NoError(t, err)
	assert.Equal(t, len(content), n)
	assert.Equal(t, content, rw.Body.Bytes())
	assert.Equal(t, http.StatusOK, scrw.statusCode)
}
