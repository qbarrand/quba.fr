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
		"../../data/images/lhc_1.jpg",
		"../../data/images/dubai_1.jpg",
		"../../data/images/singapore_1.jpg",
		"../../data/images/zermatt_1.jpg",
	}

	dimensions := []struct{ w, h int }{
		{720, 0},
		{1920, 0},
		{0, 480},
	}

	for _, p := range processors {
		p.processor.Init()
		defer p.processor.Destroy()
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
						// Without destroying vips images, vips.Shutdown() panics
						defer h.Destroy()

						h.Resize(context.Background(), 1690, 0)
						h.SetFormat(Webp)
						h.StripMetadata()

						res, err := h.Bytes()
						require.NoError(b, err)
						assert.NotEmpty(b, res)
					}
				})
			}
		}
	}
}
