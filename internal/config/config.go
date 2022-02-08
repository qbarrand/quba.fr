package config

import (
	"flag"

	"github.com/qbarrand/quba.fr/internal/handlers/sitemap"
)

type Config struct {
	Addr        string
	ImagesDir   string
	LastMod     string
	LogLevel    string
	MetricsAddr string
	WebrootDir  string
}

func ParseCommandLine(args []string) (*Config, error) {
	cfg := Config{}

	flagSet := flag.NewFlagSet("quba.fr", flag.ContinueOnError)

	flagSet.StringVar(&cfg.Addr, "addr", ":8080", "the address to listen on")
	flagSet.StringVar(&cfg.ImagesDir, "img-dir", "data/img-src", "path to the directory containing background img-src")
	flagSet.StringVar(&cfg.LastMod, "lastmod", sitemap.LastModNow(), "the last modification time of this site's contents")
	flagSet.StringVar(&cfg.LogLevel, "log-level", "info", "the log level")
	flagSet.StringVar(&cfg.MetricsAddr, "metrics-addr", ":9090", "the metrics address to listen on")
	flagSet.StringVar(&cfg.WebrootDir, "webroot-dir", "webroot", "path to the directory containing static file")

	return &cfg, flagSet.Parse(args)
}
