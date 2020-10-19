package image

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcceptHeaderToFormat(t *testing.T) {
	const (
		mimeJpeg = "image/jpeg"
		mimeWebp = "image/webp"
	)

	cases := []struct {
		types            []string
		expectedFormat   Format
		expectedMIMEType string
		returnedError    error
	}{
		{
			types:            []string{mimeJpeg},
			expectedFormat:   JPEG,
			expectedMIMEType: mimeJpeg,
		},
		{
			types:            []string{mimeWebp},
			expectedFormat:   Webp,
			expectedMIMEType: mimeWebp,
		},
		{
			types:            []string{mimeWebp, mimeJpeg},
			expectedFormat:   Webp,
			expectedMIMEType: mimeWebp,
		},
		{
			types:            []string{"image/bmp"},
			expectedFormat:   0,
			expectedMIMEType: "",
			returnedError:    ErrNotAcceptable,
		},
	}

	for _, c := range cases {
		t.Run(strings.Join(c.types, "_"), func(t *testing.T) {
			f, mime, err := AcceptHeaderToFormat(c.types)
			assert.Equal(t, c.expectedFormat, f)
			assert.Equal(t, c.expectedMIMEType, mime)

			if c.returnedError != nil {
				assert.True(t, errors.Is(ErrNotAcceptable, err))
			}
		})
	}
}
