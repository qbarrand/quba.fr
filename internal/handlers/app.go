package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	img "github.com/qbarrand/quba.fr/internal/image"
)

type AppOptions struct {
	ImageProcessor img.Processor
	LastMod        time.Time
	WebRootDir     string
}

type App struct {
	file        http.Handler
	healthz     http.Handler
	image       http.Handler
	imageLister http.Handler
	sitemap     http.Handler
}

func (a *App) Router() *mux.Router {
	router := mux.NewRouter()

	subRouter := router.Methods(http.MethodGet).Subrouter()
	subRouter.Handle("/healthz", a.healthz)
	subRouter.Handle("/sitemap.xml", a.sitemap)
	subRouter.Path("/images").Handler(a.imageLister)
	subRouter.PathPrefix("/images/").Handler(a.image)
	subRouter.PathPrefix("/").Handler(a.file)

	return router
}

func NewApp(opts *AppOptions, logger logrus.FieldLogger) (*App, error) {
	const handlerKey = "handler"

	sitemap, err := newSitemap(opts.LastMod, logger.WithField(handlerKey, "sitemap"))
	if err != nil {
		return nil, fmt.Errorf("could not create the sitemap handler: %v", err)
	}

	const subdir = "images"

	imagesPath := filepath.Join(opts.WebRootDir, subdir)

	image, err := newImage(opts.ImageProcessor, imagesPath, logger.WithField(handlerKey, "image"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the image handler: %v", err)
	}

	imageLister, err := newImageLister(&img.StaticLister{}, logger.WithField(handlerKey, "lister"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the lister handler: %w", err)
	}

	app := App{
		file:        http.FileServer(http.Dir(opts.WebRootDir)),
		healthz:     newHealthz(logger),
		image:       image,
		imageLister: imageLister,
		sitemap:     sitemap,
	}

	return &app, nil
}
