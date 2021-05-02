package sitemap

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("empty lastmod", func(t *testing.T) {
		_, err := New("")
		assert.Error(t, err)
	})

	t.Run("invalid lastmod", func(t *testing.T) {
		_, err := New("abcd")
		assert.Error(t, err)
	})

	t.Run("valid lastmod", func(t *testing.T) {
		handler, err := New("2021-05-01")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, "application/xml", w.Result().Header.Get("Content-Type"))

		const expectedBody = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>https://quba.fr/</loc>
		<lastmod>2021-05-01</lastmod>
		<changefreq>monthly</changefreq>
		<priority>1.0</priority>
	</url>
</urlset>`

		responseBody, err := io.ReadAll(w.Result().Body)
		require.NoError(t, err)

		assert.Equal(t, []byte(expectedBody), responseBody)
	})
}
