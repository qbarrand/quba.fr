package main

import (
	_ "embed"
	"errors"
	"flag"
	"net/http"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/internal/config"
	"github.com/qbarrand/quba.fr/internal/handlers"
	"github.com/qbarrand/quba.fr/pkg/httputils"
)

var version string

func main() {
	logger := logrus.New()

	cfg, err := config.ParseCommandLine(os.Args[0], os.Args[1:])
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

	opts := &handlers.AppOptions{
		ImgOutDir:  cfg.ImgOutDir,
		LastMod:    "2022-02-01",
		WebrootDir: cfg.WebrootDir,
	}

	logger.
		WithFields(logrus.Fields{
			"img-out-dir": opts.ImgOutDir,
			"last-mod":    opts.LastMod,
			"webroot-dir": opts.WebrootDir,
		}).
		Debug("Using the following configuration")

	app, err := handlers.NewApp(opts, logger)
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
		WithField("version", strings.TrimSuffix(version, "\n")).
		Info("Starting the app")

	chanErr := make(chan error)

	go func() {
		logger.
			WithField("addr", cfg.Addr).
			Info("Starting the main server")

		chanErr <- http.ListenAndServe(cfg.Addr, main)
	}()

	go func() {
		logger.
			WithField("addr", cfg.MetricsAddr).
			Info("Starting the metrics server")

		chanErr <- runMetrics(cfg.MetricsAddr)
	}()

	// go func() {
	// 	mux := http.NewServeMux()

	// 	mux.HandleFunc("/debug/pprof/", pprof.Index)
	// 	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	// 	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	// 	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	// 	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// 	chanErr <- http.ListenAndServe("localhost:6060", mux)
	// }()

	logger.WithError(<-chanErr).Fatal("General error caught")
}

func runMetrics(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, mux)
}
