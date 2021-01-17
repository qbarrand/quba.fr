package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	img "github.com/qbarrand/quba.fr/internal/image"
	"github.com/qbarrand/quba.fr/pkg/assets"
)

type AppOptions struct {
	ImageProcessor img.Processor
	LastMod        time.Time
	WebRootDir     string
}

type App struct {
	assets      http.Handler
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
	subRouter.PathPrefix("/").Handler(a.assets)

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

	meta := img.NewStaticMetaDB()

	image, err := newImage(opts.ImageProcessor, imagesPath, meta, logger.WithField(handlerKey, "image"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the image handler: %v", err)
	}

	imageLister, err := newImageLister(meta, logger.WithField(handlerKey, "lister"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the lister handler: %w", err)
	}

	assetsServer, err := assetsHandler(opts.WebRootDir+"/", logger)
	if err != nil {
		return nil, fmt.Errorf("could not build the assets handler: %v", err)
	}

	app := App{
		assets:      assetsServer.Handler(),
		healthz:     newHealthz(logger),
		image:       image,
		imageLister: imageLister,
		sitemap:     sitemap,
	}

	return &app, nil
}

func assetsHandler(staticDir string, logger logrus.FieldLogger) (*assets.Server, error) {
	s := assets.NewServer(logger)

	err := filepath.Walk(staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		refPath, err := filepath.Rel(staticDir, path)
		if err != nil {
			return fmt.Errorf("could not determine the relative path of %q from %q: %v", path, staticDir, err)
		}

		if info.IsDir() {
			if refPath == "images" {
				return filepath.SkipDir
			}

			return nil
		}

		target, err := s.AddStaticFile(path, "/"+refPath, true)
		if err != nil {
			return fmt.Errorf("could not add %q: %v", path, err)
		}

		logger.
			WithFields(logrus.Fields{
				"filepath": path,
				"refPath":  refPath,
				"target":   target,
			}).
			Debug("Added a static file")

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not add static files: %v", err)
	}

	for _, css := range []string{"/css/fontawesome-all.min.css", "/css/main.css"} {
		contents, err := ioutil.ReadFile("templates" + css)
		if err != nil {
			return nil, fmt.Errorf("could not read %s: %v", css, err)
		}

		if _, err := s.AddTemplate(string(contents), css, true); err != nil {
			return nil, fmt.Errorf("could not compile %s: %v", css, err)
		}
	}

	indexTxt, err := ioutil.ReadFile("templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("could not read the index: %v", err)
	}

	if _, err := s.AddTemplate(string(indexTxt), "/index.html", false); err != nil {
		return nil, fmt.Errorf("could not compile the index: %v", err)
	}

	if err := s.AddDirectory("/"); err != nil {
		return nil, fmt.Errorf("could not add the root directory: %v", err)
	}

	return s, nil
}
