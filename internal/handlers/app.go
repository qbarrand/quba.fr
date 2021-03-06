package handlers

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/data/images"
	"github.com/qbarrand/quba.fr/data/webroot"
	"github.com/qbarrand/quba.fr/internal/handlers/healthz"
	"github.com/qbarrand/quba.fr/internal/handlers/image"
	"github.com/qbarrand/quba.fr/internal/handlers/sitemap"
	"github.com/qbarrand/quba.fr/internal/imgpro"
)

type AppOptions struct {
	ImageProcessor imgpro.Processor
	LastMod        string
}

type App struct {
	healthz     *healthz.Healthz
	image       *image.Image
	imageLister http.HandlerFunc
	sitemap     http.HandlerFunc
	webRootFS   fs.FS
}

func (a *App) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.FS(a.webRootFS)))
	mux.Handle("/healthz", a.healthz)
	mux.HandleFunc("/images", a.imageLister)
	mux.Handle("/images/", http.StripPrefix("/images/", a.image))
	mux.HandleFunc("/sitemap.xml", a.sitemap)

	return mux
}

func NewApp(opts *AppOptions, logger logrus.FieldLogger) (*App, error) {
	const handlerKey = "handler"

	sm, err := sitemap.New(opts.LastMod)
	if err != nil {
		return nil, fmt.Errorf("could not create the sitemap handler: %v", err)
	}

	localImages := images.LocalImagesWithMetadata()

	imageHandler, err := image.New(opts.ImageProcessor, localImages, logger.WithField(handlerKey, "image"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the image handler: %v", err)
	}

	imageLister, err := image.Lister(localImages, logger.WithField(handlerKey, "lister"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the lister handler: %w", err)
	}

	app := App{
		webRootFS:   webroot.WebRoot,
		healthz:     healthz.New(logger),
		image:       imageHandler,
		imageLister: imageLister,
		sitemap:     sm,
	}

	return &app, nil
}
