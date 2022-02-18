package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCommandLine(t *testing.T) {
	const (
		addr        = "some-host:1234"
		imgOutDir   = "img-out-dir"
		logLevel    = "warn"
		metricsAddr = "some-other-host:9090"
		webrootDir  = "some-webroot-dir"
	)

	expected := &Config{
		Addr:        addr,
		ImgOutDir:   imgOutDir,
		LogLevel:    logLevel,
		MetricsAddr: metricsAddr,
		WebrootDir:  webrootDir,
	}

	args := []string{
		"--addr", addr,
		"--img-out-dir", imgOutDir,
		"--log-level", logLevel,
		"--metrics-addr", metricsAddr,
		"--webroot-dir", webrootDir,
	}

	cfg, err := ParseCommandLine("", args)

	require.NoError(t, err)
	assert.Equal(t, expected, cfg)
}
