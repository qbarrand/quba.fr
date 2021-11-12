package sitemap

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const lastModLayout = "2006-01-02"

// New returns a new handler that writes the sitemap.
func New(lastmod string, logger logrus.FieldLogger) (http.HandlerFunc, error) {
	if _, err := time.Parse(lastModLayout, lastmod); err != nil {
		return nil, fmt.Errorf("invalid lastmod: %v", err)
	}

	const sitemapTemplateStr = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>https://quba.fr/</loc>
		<lastmod>%s</lastmod>
		<changefreq>monthly</changefreq>
		<priority>1.0</priority>
	</url>
</urlset>`

	body := []byte(fmt.Sprintf(sitemapTemplateStr, lastmod))

	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(body); err != nil {
			logger.WithError(err).Error("Could not write the body")
		}
	}

	return handler, nil
}

func LastModNow() string {
	return time.Now().Format(lastModLayout)
}
