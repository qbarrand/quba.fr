package httplog

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoggingResponseWriter(t *testing.T) {
	t.Parallel()

	rw := httptest.NewRecorder()
	lrw := NewLoggingResponseWriter(rw)

	assert.Equal(t, rw, lrw.wrapped)
	assert.Implements(t, (*http.ResponseWriter)(nil), lrw)
}

func TestLoggingResponseWriter_Header(t *testing.T) {
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

	lrw := LoggingResponseWriter{wrapped: rw}

	assert.Equal(t, headers, lrw.Header())
}

func TestLoggingResponseWriter_StatusCode(t *testing.T) {
	t.Parallel()

	const status = http.StatusTeapot

	lrw := LoggingResponseWriter{statusCode: status}

	assert.Equal(t, status, lrw.StatusCode())
}

func TestLoggingResponseWriter_Write(t *testing.T) {
	t.Parallel()

	content := []byte("some-content")

	rw := httptest.NewRecorder()
	lrw := &LoggingResponseWriter{wrapped: rw}

	n, err := lrw.Write(content)

	require.NoError(t, err)
	assert.Equal(t, len(content), n)

	responseBody, err := ioutil.ReadAll(rw.Result().Body)

	require.NoError(t, err)
	assert.Equal(t, content, responseBody)
}
