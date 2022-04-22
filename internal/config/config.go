package config

import (
	"flag"
)

type Config struct {
	Addr        string
	ImgOutDir   string
	LogLevel    string
	MetricsAddr string
	WebrootDir  string
}

func ParseCommandLine(programName string, args []string) (*Config, error) {
	cfg := Config{}

	flagSet := flag.NewFlagSet(programName, flag.ContinueOnError)

	flagSet.StringVar(&cfg.Addr, "addr", ":8080", "the address to listen on")
	flagSet.StringVar(&cfg.ImgOutDir, "img-out-dir", "img-out", "path to the directory containing background images")
	flagSet.StringVar(&cfg.LogLevel, "log-level", "info", "the log level")
	flagSet.StringVar(&cfg.MetricsAddr, "metrics-addr", ":9090", "the metrics address to listen on")
	flagSet.StringVar(&cfg.WebrootDir, "webroot-dir", "webroot", "path to the directory containing static file")

	return &cfg, flagSet.Parse(args)
}
