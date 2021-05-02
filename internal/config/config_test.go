package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_configFromArgs(t *testing.T) {
	const (
		addr     = "some-host:1234"
		lastMod  = "2001-02-03"
		logLevel = "warn"
	)

	expected := &Config{
		Addr:     addr,
		LastMod:  lastMod,
		LogLevel: logLevel,
	}

	args := []string{
		"--addr", addr,
		"--lastmod", lastMod,
		"--log-level", logLevel,
	}

	cfg, err := ParseCommandLine(args)

	require.NoError(t, err)
	assert.Equal(t, expected, cfg)
}
