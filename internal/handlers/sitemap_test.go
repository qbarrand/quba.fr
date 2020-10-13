package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newSitemap(t *testing.T) {
	t.Run("empty lastmod", func(t *testing.T) {
		_, err := newSitemap("")
		assert.Error(t, err)
	})

	t.Run("should work", func(t *testing.T) {
		s, err := newSitemap("2020-10-14")
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})
}

func TestSitemap_ServeHTTP(t *testing.T) {
	s, err := newSitemap("2020-10-14")

	require.NoError(t, err)

	const expected = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>https://quba.fr/</loc>
		<lastmod>2020-10-14</lastmod>
		<changefreq>monthly</changefreq>
		<priority>1.0</priority>
	</url>
</urlset>`

	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, nil)

	assert.Equal(t, expected, recorder.Body.String())
}
