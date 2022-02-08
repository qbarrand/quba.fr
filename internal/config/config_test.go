package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommandLine(t *testing.T) {
	const (
		addr        = "some-host:1234"
		imagesDir   = "some-img-src-dir"
		lastMod     = "2001-02-03"
		logLevel    = "warn"
		metricsAddr = "some-other-host:9090"
		webrootDir  = "some-webroot-dir"
	)

	expected := &Config{
		Addr:        addr,
		LastMod:     lastMod,
		LogLevel:    logLevel,
		MetricsAddr: metricsAddr,
	}

	args := []string{
		"--addr", addr,
		"--img-src-dir", imagesDir,
		"--lastmod", lastMod,
		"--log-level", logLevel,
		"--metrics-addr", metricsAddr,
		"--webroot-dir", webrootDir,
	}

	cfg, err := ParseCommandLine(args)

	require.NoError(t, err)
	assert.Equal(t, expected, cfg)
}
