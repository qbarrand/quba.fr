package handlers

import (
	"fmt"
	"net/http"

	"github.com/qbarrand/quba.fr/internal/handlers/background"
	"github.com/qbarrand/quba.fr/internal/handlers/healthz"
	"github.com/qbarrand/quba.fr/internal/handlers/sitemap"
	"github.com/sirupsen/logrus"
)

const bgImagesPath = "/images/bg/"

type AppOptions struct {
	ImgOutDir  string
	LastMod    string
	WebrootDir string
}

type App struct {
	healthz     *healthz.Healthz
	bg          *background.Handler
	bgStatic    http.Handler
	sitemap     http.HandlerFunc
	rootHandler http.Handler
}

func (a *App) Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/healthz", a.healthz)
	mux.Handle("/background", a.bg)
	mux.Handle(bgImagesPath, a.bgStatic)
	mux.Handle("/", a.rootHandler)

	return mux
}

func NewApp(opts *AppOptions, logger logrus.FieldLogger) (*App, error) {
	sm, err := sitemap.New(opts.LastMod, logger.WithField("handler", "sitemap"))
	if err != nil {
		return nil, fmt.Errorf("could not create the sitemap handler: %v", err)
	}

	bg, err := background.NewHandler(opts.ImgOutDir, bgImagesPath, logger.WithField("handler", "background"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the background handler: %v", err)
	}

	app := App{
		bg: bg,
		bgStatic: http.StripPrefix(
			bgImagesPath,
			http.FileServer(
				http.Dir(opts.ImgOutDir),
			),
		),
		healthz: healthz.New(logger),
		sitemap: sm,
		rootHandler: http.FileServer(
			http.Dir(opts.WebrootDir),
		),
	}

	return &app, nil
}
