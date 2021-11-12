package imgpro

import (
	"context"
	"fmt"

	"github.com/davidbyttow/govips/v2/vips"
)

type VipsHandler struct {
	format        Format
	ref           *vips.ImageRef
	stripMetadata bool
}

func (vh *VipsHandler) Bytes() ([]byte, error) {
	switch vh.format {
	case Webp:
		p := vips.NewWebpExportParams()
		p.StripMetadata = vh.stripMetadata
		buf, _, err := vh.ref.ExportWebp(p)
		return buf, err
	case JPEG:
		p := vips.NewJpegExportParams()
		p.StripMetadata = vh.stripMetadata
		buf, _, err := vh.ref.ExportJpeg(p)
		return buf, err
	}

	return nil, fmt.Errorf("%s: unhandled format", vh.format)
}

func (vh *VipsHandler) Destroy() error {
	vh.ref.Close()

	return nil
}

func (vh *VipsHandler) Resize(ctx context.Context, w, h int) error {
	var scale float64 = 1

	if w != 0 {
		scale = float64(w) / float64(vh.ref.Width())
	} else if h != 0 {
		scale = float64(h) / float64(vh.ref.Height())
	}

	if scale > 1 {
		return nil
	}

	chanErr := make(chan error, 1)

	go func() {
		chanErr <- vh.ref.Resize(scale, vips.KernelAuto) // vh.ref.Thumbnail(w, h, vips.InterestingAll)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-chanErr:
		return err
	}
}

func (vh *VipsHandler) SetFormat(format Format) error {
	vh.format = format
	return nil
}

func (vh *VipsHandler) StripMetadata() error {
	vh.stripMetadata = true
	return nil
}

type VipsProcessor struct{}

func (vp *VipsProcessor) BestFormats() []Format {
	return []Format{Webp, JPEG}
}

func (vp *VipsProcessor) Destroy() error {
	vips.Shutdown()

	return nil
}

func (vp *VipsProcessor) Init() error {
	vips.Startup(nil)

	return nil
}

func (vp *VipsProcessor) HandlerFromBytes(b []byte) (Handler, error) {
	ref, err := vips.NewImageFromBuffer(b)
	if err != nil {
		return nil, fmt.Errorf("could not create the handler: %v", err)
	}

	return &VipsHandler{ref: ref}, nil
}

func (vp *VipsProcessor) String() string {
	return "Vips"
}
