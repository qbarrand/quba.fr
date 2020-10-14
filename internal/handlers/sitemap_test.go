package handlers

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newSitemap(t *testing.T) {
	logger, _ := test.NewNullLogger()

	t.Run("should work", func(t *testing.T) {
		s, err := newSitemap(time.Now(), logger)
		assert.NoError(t, err)
		assert.NotNil(t, s)
	})
}

func TestSitemap_ServeHTTP(t *testing.T) {
	logger, _ := test.NewNullLogger()

	lastModTime := time.Date(2020, time.October, 14, 0, 0, 0, 0, time.UTC)

	s, err := newSitemap(lastModTime, logger)

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
