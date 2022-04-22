// +build magick

package imgpro

import (
	"context"
	"fmt"

	"gopkg.in/gographics/imagick.v3/imagick"
)

type ImageMagickHandler struct {
	mw *imagick.MagickWand
}

func (imh *ImageMagickHandler) Bytes() ([]byte, error) {
	return imh.mw.GetImageBlob(), nil
}

func (imh *ImageMagickHandler) Destroy() error {
	imh.mw.Destroy()

	return nil
}

func (imh *ImageMagickHandler) Resize(ctx context.Context, w, h int) error {
	uw := uint(w)
	uh := uint(h)

	if w == 0 || h == 0 {
		ow := imh.mw.GetImageWidth()
		oh := imh.mw.GetImageHeight()

		owf := float64(ow)
		ohf := float64(oh)

		if w == 0 {
			if scale := ohf / float64(h); scale < 1 {
				uw = uint(scale * owf)
			} else {
				uw = uint(ow)
				uh = uint(oh)
			}
		} else if h == 0 {
			if scale := owf / float64(w); scale < 1 {
				uh = uint(scale * ohf)
			} else {
				uw = uint(ow)
				uh = uint(oh)
			}
		}
	}

	c := make(chan error)

	go func() {
		c <- imh.mw.AdaptiveResizeImage(uw, uh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func (imh *ImageMagickHandler) SetFormat(f Format) error {
	switch f {
	case Webp:
		return imh.mw.SetFormat("WEBP")
	case JPEG:
		return imh.mw.SetFormat("JPEG")
	}

	return fmt.Errorf("%s: no corresponding format found for ImageMagick", f)
}

func (imh *ImageMagickHandler) StripMetadata() error {
	return imh.mw.StripImage()
}

type ImageMagickProcessor struct{}

func (imp *ImageMagickProcessor) BestFormats() []Format {
	return []Format{Webp, JPEG}
}

func (imp *ImageMagickProcessor) Destroy() error {
	imagick.Terminate()

	return nil
}

func (imp *ImageMagickProcessor) Init(_ int) error {
	imagick.Initialize()

	return nil
}

func (imp *ImageMagickProcessor) HandlerFromBytes(b []byte) (Handler, error) {
	imh := ImageMagickHandler{
		mw: imagick.NewMagickWand(),
	}

	return &imh, imh.mw.ReadImageBlob(b)
}

func (imp *ImageMagickProcessor) String() string {
	return "ImageMagick"
}
