package main

import (
	"flag"
	"time"
)

type config struct {
	addr     string
	lastMod  time.Time
	logLevel string
}

func configFromArgs(args []string) (*config, error) {
	cfg := config{}

	flagSet := flag.NewFlagSet("quba.fr", flag.ContinueOnError)

	flagSet.StringVar(&cfg.addr, "addr", ":8080", "the address to listen on")
	flagSet.StringVar(&cfg.logLevel, "log-level", "info", "the log level")

	return &cfg, flagSet.Parse(args)
}
