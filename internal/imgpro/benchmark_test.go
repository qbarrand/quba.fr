//go:build magick
// +build magick

package imgpro

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkSuite(b *testing.B) {
	vips.LoggingSettings(func(messageDomain string, messageLevel vips.LogLevel, message string) {}, 0)

	processors := []struct {
		name      string
		processor Processor
	}{
		{
			name:      "ImageMagick",
			processor: &ImageMagickProcessor{},
		},
		{
			name:      "vips",
			processor: &VipsProcessor{},
		},
	}

	images := []string{
		"../../img-src/lhc_1.jpg",
		"../../img-src/dubai_1.jpg",
		"../../img-src/singapore_1.jpg",
		"../../img-src/zermatt_1.jpg",
	}

	dimensions := []struct{ w, h int }{
		{720, 0},
		{1920, 0},
		{0, 480},
	}

	for _, p := range processors {
		err := p.processor.Init(1)
		require.NoError(b, err)

		defer func(proc Processor) {
			require.NoError(
				b,
				proc.Destroy(),
			)
		}(p.processor)
	}

	for _, img := range images {
		imageBytes, err := os.ReadFile(img)
		require.NoError(b, err)

		imgShort := filepath.Base(img)

		for _, d := range dimensions {
			for _, p := range processors {
				b.Run(fmt.Sprintf("%s_%dx%d_%s", imgShort, d.w, d.h, p.name), func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						h, err := p.processor.HandlerFromBytes(imageBytes)
						require.NoError(b, err)

						// Without destroying vips img-src, vips.Shutdown() panics
						defer func() {
							err = h.Destroy()
							require.NoError(b, err)
						}()

						err = h.Resize(context.Background(), 1690, 0)
						require.NoError(b, err)

						err = h.SetFormat(Webp)
						require.NoError(b, err)

						err = h.StripMetadata()
						require.NoError(b, err)

						res, err := h.Bytes()
						require.NoError(b, err)
						assert.NotEmpty(b, res)
					}
				})
			}
		}
	}
}
