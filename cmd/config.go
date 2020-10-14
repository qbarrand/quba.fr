package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/qbarrand/quba.fr/internal/handlers"
)

var errBadLastMod = errors.New("could not parse lastmod")

type config struct {
	addr     string
	lastMod  time.Time
	logLevel string
}

func configFromArgs(args []string) (*config, error) {
	cfg := config{}

	var lastMod string

	flagSet := flag.NewFlagSet("quba.fr", flag.ContinueOnError)

	flagSet.StringVar(&cfg.addr, "addr", ":8080", "the address to listen on")
	flagSet.StringVar(&lastMod, "lastmod", handlers.TimeToLastMod(time.Now()), "the last modification time of this site's contents")
	flagSet.StringVar(&cfg.logLevel, "log-level", "info", "the log level")

	var err error

	if err = flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf("could not parse the arguments: %w", err)
	}

	cfg.lastMod, err = handlers.TimeFromLastMod(lastMod)
	if err != nil {
		return nil, fmt.Errorf("%w: could not parse the last modification time: %v", errBadLastMod, err)
	}

	return &cfg, nil
}
