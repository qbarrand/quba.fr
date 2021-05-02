package config

import (
	"flag"

	"github.com/qbarrand/quba.fr/internal/handlers/sitemap"
)

type Config struct {
	Addr     string
	LastMod  string
	LogLevel string
}

func ParseCommandLine(args []string) (*Config, error) {
	cfg := Config{}

	flagSet := flag.NewFlagSet("quba.fr", flag.ContinueOnError)

	flagSet.StringVar(&cfg.Addr, "addr", ":8080", "the address to listen on")
	flagSet.StringVar(&cfg.LastMod, "lastmod", sitemap.LastModNow(), "the last modification time of this site's contents")
	flagSet.StringVar(&cfg.LogLevel, "log-level", "info", "the log level")

	return &cfg, flagSet.Parse(args)
}
