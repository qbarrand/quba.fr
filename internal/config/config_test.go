package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommandLine(t *testing.T) {
	const (
		addr           = "some-host:1234"
		imageProcessor = "vips"
		lastMod        = "2001-02-03"
		logLevel       = "warn"
		metricsAddr    = "some-other-host:9090"
	)

	expected := &Config{
		Addr:           addr,
		ImageProcessor: imageProcessor,
		LastMod:        lastMod,
		LogLevel:       logLevel,
		MetricsAddr:    metricsAddr,
	}

	args := []string{
		"--addr", addr,
		"--image-processor", imageProcessor,
		"--lastmod", lastMod,
		"--log-level", logLevel,
		"--metrics-addr", metricsAddr,
	}

	cfg, err := ParseCommandLine(args)

	require.NoError(t, err)
	assert.Equal(t, expected, cfg)
}
