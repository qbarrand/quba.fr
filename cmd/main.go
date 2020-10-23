package main

import (
	"errors"
	"flag"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/internal/handlers"
	"github.com/qbarrand/quba.fr/internal/image"
	"github.com/qbarrand/quba.fr/pkg/httputils"
)

func main() {
	logger := logrus.New()

	cfg, err := configFromArgs(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}

		logger.WithError(err).Fatal("Could not parse the command line")
	}

	logLevel, err := logrus.ParseLevel(cfg.logLevel)
	if err != nil {
		logger.WithError(err).Fatal("Could not parse the log level")
	}

	logger.SetLevel(logLevel)

	opts := handlers.AppOptions{
		ImageProcessor: &image.VipsProcessor{},
		LastMod:        cfg.lastMod,
		WebRootDir:     "webroot",
	}

	app, err := handlers.NewApp(&opts, logger)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize the app")
	}

	// Intercept all requests, then forward them to the router.
	// We use this instead of a middleware, as those are only hit when the router
	// has a match (not for 404, 405 etc).
	main := httputils.LoggingMiddleware(
		logger,
		app.Router(),
	)

	logger.WithField("addr", cfg.addr).Info("Starting the server")

	if err := http.ListenAndServe(cfg.addr, main); err != nil {
		logger.WithError(err).Fatal("General error caught")
	}
}
