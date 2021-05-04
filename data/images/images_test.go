package images

import (
	"image/jpeg"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalImagesWithMetadata(t *testing.T) {
	entries, err := filepath.Glob("*.jpg")
	require.NoError(t, err)

	assert.NotEmpty(t, entries)

	// iwn contains an embed.FS that should contain all files in the current directory
	iwm := LocalImagesWithMetadata()

	assert.NoError(
		t,
		fstest.TestFS(iwm, entries...),
	)

	for _, e := range entries {
		fd, meta, err := iwm.OpenWithMetadata(e)

		assert.NotNil(t, fd)
		assert.NotNil(t, meta)
		require.NoError(t, err)
	}
}

func TestEmbedded_OpenWithMetadata(t *testing.T) {
	em := &embedded{fs: local}

	fd, meta, err := em.OpenWithMetadata("dents_du_midi_1.jpg")
	require.NoError(t, err)

	// Check the image is indeed jpeg
	_, err = jpeg.Decode(fd)
	assert.NoError(t, err)

	// Check the metadata is correct
	expectedMeta := &Metadata{
		Date:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		Location:  "Dents du Midi, Switzerland",
		MainColor: "#4279AC",
	}

	assert.Equal(t, expectedMeta, meta)
}
