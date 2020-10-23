package httputils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequestID(t *testing.T) {
	t.Run("no request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/some/url", nil)
		assert.Equal(t, "<nil>", GetRequestID(req))
	})

	t.Run("struct request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/some/url", nil)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				requestIDKey,
				struct{ Value int }{Value: 12345},
			),
		)

		assert.Equal(t, "<error>", GetRequestID(req))
	})

	t.Run("integer request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/some/url", nil)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				requestIDKey,
				12345,
			),
		)

		assert.Equal(t, "12345", GetRequestID(req))
	})

	t.Run("string request ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/some/url", nil)
		req = req.WithContext(
			context.WithValue(
				req.Context(),
				requestIDKey,
				"12345",
			),
		)

		assert.Equal(t, "12345", GetRequestID(req))
	})
}
