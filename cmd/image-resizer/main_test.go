package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestImageDir(t *testing.T) {
	getMetadata := func(t *testing.T) Metadata {
		t.Helper()

		md, err := ReadFromFile("../../img-src/metadata.json")
		require.NoError(t, err)

		return md
	}

	t.Run("all images should have a manifest", func(t *testing.T) {
		md := getMetadata(t)

		files, err := filepath.Glob("../../img-src/*.jpg")
		require.NoError(t, err)

		for _, f := range files {
			require.Contains(t, md, filepath.Base(f))
		}
	})

	t.Run("all manifests should target an image", func(t *testing.T) {

		for name := range getMetadata(t) {
			fd, err := os.Open("../../img-src/" + name)
			require.NoError(t, err)

			err = fd.Close()
			require.NoError(t, err)
		}
	})
}
