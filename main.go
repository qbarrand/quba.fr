package main

import (
	_ "embed"
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/internal/config"
	"github.com/qbarrand/quba.fr/internal/handlers"
	"github.com/qbarrand/quba.fr/internal/imgpro"
	"github.com/qbarrand/quba.fr/pkg/httputils"
)

//go:embed VERSION
var version string

func main() {
	logger := logrus.New()

	cfg, err := config.ParseCommandLine(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}

		logger.WithError(err).Fatal("Could not parse the command line")
	}

	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logger.WithError(err).Fatal("Could not parse the log level")
	}

	logger.SetLevel(logLevel)

	opts := handlers.AppOptions{
		ImageProcessor: &imgpro.VipsProcessor{},
		LastMod:        cfg.LastMod,
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

	logger.
		WithFields(logrus.Fields{
			"addr":    cfg.Addr,
			"version": strings.TrimSuffix(version, "\n"),
		}).
		Info("Starting the server")

	if err := http.ListenAndServe(cfg.Addr, main); err != nil {
		logger.WithError(err).Fatal("General error caught")
	}
}
