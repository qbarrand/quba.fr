package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/davidbyttow/govips/pkg/vips"
)

var (
	ErrFormatUnavailable = errors.New("this format is not available in vips")

	vipsFormats = []Format{Webp, JPEG}
)

type VipsHandler struct {
	export *vips.ExportParams
	ref    *vips.ImageRef
}

func (vh *VipsHandler) Bytes() ([]byte, error) {
	buf, _, err := vh.ref.Export(*vh.export)

	vh.ref.Close()

	return buf, err
}

func (vh *VipsHandler) Destroy() error {
	vh.ref.Close()

	return nil
}

func (vh *VipsHandler) Resize(ctx context.Context, w, h int) error {
	chanErr := make(chan error, 1)

	go func() {
		chanErr <- vh.ref.ThumbnailImage(w)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chanErr:
		return err
	}
}

func (vh *VipsHandler) SetFormat(format Format) error {
	var err error

	vh.export.Format, err = formatToVipsImageType(format)

	return err
}

type VipsProcessor struct{}

func (vp *VipsProcessor) BestFormats() []Format {
	return vipsFormats
}

func (vp *VipsProcessor) Destroy() error {
	vips.Shutdown()

	return nil
}

func (vp *VipsProcessor) Init() error {
	vips.Startup(nil)

	return nil
}

func (vp *VipsProcessor) NewImageHandler(s string) (Handler, error) {
	ref, err := vips.NewImageFromFile(s)
	if err != nil {
		return nil, fmt.Errorf("could not create the handler: %v", err)
	}

	vh := VipsHandler{
		export: &vips.ExportParams{},
		ref:    ref,
	}

	return &vh, nil
}

func (vp *VipsProcessor) HandlerFromBytes(b []byte) (Handler, error) {
	ref, err := vips.NewImageFromBuffer(b)
	if err != nil {
		return nil, fmt.Errorf("could not create the handler: %v", err)
	}

	vh := VipsHandler{
		export: &vips.ExportParams{},
		ref:    ref,
	}

	return &vh, nil
}

func formatToVipsImageType(format Format) (vips.ImageType, error) {
	it := vips.ImageTypeUnknown

	switch format {
	case JPEG:
		it = vips.ImageTypeJPEG
	case Webp:
		it = vips.ImageTypeWEBP
	}

	if it == 0 {
		return it, ErrFormatUnavailable
	}

	return it, nil
}
