package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_configFromArgs(t *testing.T) {
	t.Run("should work as expected", func(t *testing.T) {
		const (
			addr     = "some-host:1234"
			lastMod  = "2001-02-03"
			logLevel = "warn"
		)

		expected := &config{
			addr:     addr,
			lastMod:  time.Date(2001, time.February, 3, 0, 0, 0, 0, time.UTC),
			logLevel: logLevel,
		}

		args := []string{
			"--addr", addr,
			"--lastmod", lastMod,
			"--log-level", logLevel,
		}

		cfg, err := configFromArgs(args)

		require.NoError(t, err)
		assert.Equal(t, expected, cfg)
	})

	t.Run("malformed lastmod", func(t *testing.T) {
		_, err := configFromArgs([]string{"--lastmod", "abcd"})
		require.True(t, errors.Is(err, errBadLastMod))
	})
}
