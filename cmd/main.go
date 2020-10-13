package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/internal/handlers"
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
		LastMod: time.Now(),
	}

	app, err := handlers.NewApp(&opts, logger)
	if err != nil {
		logger.WithError(err).Fatal("Could not initialize the app")
	}

	if err := http.ListenAndServe(cfg.addr, app); err != nil {
		logger.WithError(err).Fatal("General error caught")
	}
}
