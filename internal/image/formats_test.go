package image

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAcceptHeaderToFormat(t *testing.T) {
	cases := []struct {
		types         []string
		expected      Format
		returnedError error
	}{
		{
			types:    []string{"image/jpeg"},
			expected: JPEG,
		},
		{
			types:    []string{"image/webp"},
			expected: Webp,
		},
		{
			types:    []string{"image/webp", "image/jpeg"},
			expected: Webp,
		},
		{
			types:         []string{"image/bmp"},
			returnedError: ErrNotAcceptable,
		},
	}

	for _, c := range cases {
		t.Run(strings.Join(c.types, "_"), func(t *testing.T) {
			f, err := AcceptHeaderToFormat(c.types)
			assert.Equal(t, c.expected, f)

			if c.returnedError != nil {
				assert.True(t, errors.Is(ErrNotAcceptable, err))
			}
		})
	}
}
