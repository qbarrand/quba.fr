package handlers

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/internal/handlers/healthz"
	"github.com/qbarrand/quba.fr/internal/handlers/sitemap"
)

type AppOptions struct {
	ImagesDir  string
	LastMod    string
	WebrootDir string
}

type App struct {
	healthz    *healthz.Healthz
	imagesDir  string
	sitemap    http.HandlerFunc
	webRootDir string
}

func (a *App) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle(
		"/",
		http.FileServer(
			http.Dir(
				a.webRootDir,
			),
		),
	)
	mux.Handle("/healthz", a.healthz)

	mux.Handle(
		"/img-src/",
		http.StripPrefix(
			"/img-src/",
			http.FileServer(
				http.Dir(a.imagesDir),
			),
		),
	)

	mux.HandleFunc("/sitemap.xml", a.sitemap)

	return mux
}

func NewApp(opts *AppOptions, logger logrus.FieldLogger) (*App, error) {
	sm, err := sitemap.New(opts.LastMod)
	if err != nil {
		return nil, fmt.Errorf("could not create the sitemap handler: %v", err)
	}

	app := App{
		healthz:    healthz.New(logger),
		imagesDir:  opts.ImagesDir,
		sitemap:    sm,
		webRootDir: opts.WebrootDir,
	}

	return &app, nil
}
