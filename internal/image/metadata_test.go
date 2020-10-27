package image

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newMetadata(t *testing.T) {
	t.Parallel()

	const (
		month    = time.June
		year     = 2020
		location = "test-location"
	)

	details := newMetadata(month, year, location)

	expected := &Metadata{
		Date:     time.Date(year, month, 0, 0, 0, 0, 0, time.UTC),
		Location: location,
	}

	assert.Equal(t, expected, details)
}

func TestStaticMetaDB_AllNames(t *testing.T) {
	t.Parallel()

	smdb := NewStaticMetaDB()

	allNames, err := smdb.AllNames()

	require.NoError(t, err)

	t.Run("filenames from the internal map", func(t *testing.T) {
		m := make(map[string]bool, len(allNames))

		for _, name := range allNames {
			assert.NotContains(t, m, name)
			m[name] = true
		}
	})

	t.Run("all files exist", func(t *testing.T) {
		for _, name := range allNames {
			path := filepath.Join("../../webroot/images", name)

			_, err := os.Stat(path)
			assert.NoError(t, err, "Could not stat(%q)", path)
		}
	})
}

func TestStaticMetaDB_GetMetadata(t *testing.T) {
	t.Parallel()

	smdb := NewStaticMetaDB()

	t.Run("dubai_1.jpg", func(t *testing.T) {
		meta, err := smdb.GetMetadata("dubai_1.jpg")

		require.NoError(t, err)

		expectedMeta := Metadata{
			Date:     time.Date(2017, time.June, 0, 0, 0, 0, 0, time.UTC),
			Location: "Dubai, UAE",
		}

		assert.Equal(t, &expectedMeta, meta)
	})

	t.Run("non existent key", func(t *testing.T) {
		_, err := smdb.GetMetadata("non_existent_key")
		assert.True(t, errors.Is(err, ErrNoSuchMetadata))
	})
}
