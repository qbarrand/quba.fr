package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"github.com/sirupsen/logrus"
)

type sitemap struct {
	buffer bytes.Buffer
	logger logrus.FieldLogger
}

func newSitemap(lastmod string, logger logrus.FieldLogger) (*sitemap, error) {
	if lastmod == "" {
		return nil, errors.New("lastMod cannot be empty")
	}

	const sitemapTemplateStr = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
	<url>
		<loc>https://quba.fr/</loc>
		<lastmod>{{ .LastMod }}</lastmod>
		<changefreq>monthly</changefreq>
		<priority>1.0</priority>
	</url>
</urlset>`

	tmpl, err := template.New("sitemap").Parse(sitemapTemplateStr)
	if err != nil {
		return nil, fmt.Errorf("could not parse the sitemap template: %v", err)
	}

	s := sitemap{
		logger: logger,
	}

	return &s, tmpl.Execute(&s.buffer, struct{ LastMod string }{LastMod: lastmod})
}

func (s *sitemap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := s.buffer.WriteTo(w); err != nil {
		s.logger.WithError(err).Error("Error while writing the sitemap")
	}
}
