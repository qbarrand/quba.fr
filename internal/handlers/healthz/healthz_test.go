package healthz

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func Test_newHealth(t *testing.T) {
	logger, _ := test.NewNullLogger()

	h := New(logger)

	assert.NotNil(t, h)
	assert.NotNil(t, h.cache)
	assert.NotNil(t, h.dnsQueryier)
	assert.Equal(t, logger, h.logger)
}

func TestHealth_ServeHTTP(t *testing.T) {
	const path = "/healthz"

	logger, _ := test.NewNullLogger()

	cases := []struct {
		expectedCacheValue bool
		expectedCode       int
		name               string
		queryier           func(string) ([]string, error)
	}{
		{
			name: "empty DNS records",
			queryier: func(_ string) ([]string, error) {
				return nil, nil
			},
			expectedCode:       http.StatusInternalServerError,
			expectedCacheValue: false,
		},
		{
			name: "too many DNS records",
			queryier: func(_ string) ([]string, error) {
				return []string{"a", "b"}, nil
			},
			expectedCode:       http.StatusInternalServerError,
			expectedCacheValue: false,
		},
		{
			name: "DNS query error",
			queryier: func(_ string) ([]string, error) {
				return nil, errors.New("whatever")
			},
			expectedCode:       http.StatusInternalServerError,
			expectedCacheValue: false,
		},
		{
			name: "unexpected TXT contents",
			queryier: func(_ string) ([]string, error) {
				return []string{"test"}, nil
			},
			expectedCode:       http.StatusInternalServerError,
			expectedCacheValue: false,
		},
		{
			name: "expected TXT contents",
			queryier: func(q string) ([]string, error) {
				// make sure we query the right thing
				//assert.Equal(t, "ping.quba.fr", q)

				return []string{"quentin@quba.fr"}, nil
			},
			expectedCode:       http.StatusOK,
			expectedCacheValue: true,
		},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s: %d", c.name, c.expectedCode), func(t *testing.T) {
			handler := &Healthz{
				cache:       &boolCache{},
				dnsQueryier: c.queryier,
				logger:      logger,
			}

			assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, "", nil, c.expectedCode)
			assert.Equal(t, c.expectedCacheValue, handler.cache.valid)
		})
	}

	t.Run("check query", func(t *testing.T) {
		handler := &Healthz{
			cache: &boolCache{},
			dnsQueryier: func(q string) ([]string, error) {
				assert.Equal(t, "ping.quba.fr", q)

				return make([]string, 0), nil
			},
			logger: logger,
		}

		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest(http.MethodGet, path, nil),
		)
	})

	t.Run("valid cache", func(t *testing.T) {
		handler := &Healthz{
			cache: &boolCache{
				lastCheck: time.Now(),
				valid:     true,
			},
			// should never be called
			dnsQueryier: nil,
			logger:      logger,
		}

		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest(http.MethodGet, path, nil),
		)
	})

	t.Run("invalid cache", func(t *testing.T) {
		called := false

		defer func() {
			if !called {
				t.Fatal("dnsQueryier not called")
			}
		}()

		handler := &Healthz{
			cache: &boolCache{
				lastCheck: time.Now().Add(-121 * time.Second),
				valid:     true,
			},
			dnsQueryier: func(_ string) ([]string, error) {
				// this value will be tested in the deferred function
				called = true
				return make([]string, 0), nil
			},
			logger: logger,
		}

		handler.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest(http.MethodGet, path, nil),
		)
	})
}
