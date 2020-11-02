package image

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcceptHeaderToFormat(t *testing.T) {
	const (
		mimeJpeg = "image/jpeg"
		mimeWebp = "image/webp"
	)

	cases := []struct {
		expectedFormat Format
		mimeType       string
		returnedError  error
	}{
		{
			mimeType:       mimeJpeg,
			expectedFormat: JPEG,
		},
		{
			mimeType:       mimeWebp,
			expectedFormat: Webp,
		},
		{
			expectedFormat: "",
			mimeType:       "abcd",
			returnedError:  ErrNotAcceptable,
		},
		{
			expectedFormat: "",
			mimeType:       "",
			returnedError:  ErrNotAcceptable,
		},
	}

	for _, c := range cases {
		t.Run(c.mimeType, func(t *testing.T) {
			format, err := FormatFromMIMEType(c.mimeType)
			assert.Equal(t, c.expectedFormat, format)

			if c.returnedError != nil {
				assert.True(t, errors.Is(ErrNotAcceptable, err))
			}
		})
	}
}
