//go:generate go run github.com/golang/mock/mockgen -destination mock_http.go -package handlers net/http Handler

package handlers

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/qbarrand/quba.fr/internal/image/mock_image"
)

func TestNewApp(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	processor := mock_image.NewMockProcessor(ctrl)

	const webroot = "/some/path"

	opts := AppOptions{
		ImageProcessor: processor,
		LastMod:        time.Now(),
		WebRootDir:     "/some/path",
	}

	logger, _ := test.NewNullLogger()

	processor.EXPECT().Init()

	app, err := NewApp(&opts, logger)

	assert.NoError(t, err)

	assert.Equal(t, http.FileServer(http.Dir(webroot)), app.file)

	h, ok := app.healthz.(*healthz)
	assert.True(t, ok)
	assert.Equal(t, logger, h.logger)

	i, ok := app.image.(*image)
	assert.True(t, ok)
	assert.Equal(t, processor, i.processor)
	assert.Equal(t, webroot+"/images", i.path)
	assert.Equal(t, logger.WithField("handler", "image"), i.logger)

	s, ok := app.sitemap.(*sitemap)
	assert.True(t, ok)
	assert.Equal(t, logger.WithField("handler", "sitemap"), s.logger)
}

func TestApp_Router(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	file := NewMockHandler(ctrl)
	healthz := NewMockHandler(ctrl)
	image := NewMockHandler(ctrl)
	sitemap := NewMockHandler(ctrl)

	app := App{
		file:    file,
		healthz: healthz,
		image:   image,
		sitemap: sitemap,
	}

	router := app.Router()

	cases := []struct {
		mock *MockHandler
		url  string
	}{
		{
			mock: file,
			url:  "/index.html",
		},
		{
			mock: healthz,
			url:  "/healthz",
		},
		{
			mock: image,
			url:  "/images/dubai_1.jpg",
		},
		{
			mock: sitemap,
			url:  "/sitemap.xml",
		},
	}

	for _, c := range cases {
		t.Run("GET "+c.url, func(t *testing.T) {
			const status = http.StatusTeapot

			c.mock.EXPECT().ServeHTTP(gomock.Any(), gomock.Any()).Do(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			})

			assert.HTTPStatusCode(t, router.ServeHTTP, http.MethodGet, c.url, nil, status)
		})
	}
}
