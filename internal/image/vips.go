package image

import (
	"context"
	"errors"
	"io"
)

var errNotImplemented = errors.New("not implemented")

type VipsHandler struct {
}

func (vh *VipsHandler) Resize(ctx context.Context, i int, i2 int) error {
	return errNotImplemented
}

func (vh *VipsHandler) SetFormat(format Format) error {
	return errNotImplemented
}

func (vh *VipsHandler) WriteTo(w io.Writer) (n int64, err error) {
	return 0, errNotImplemented
}

type VipsProcessor struct {
}

func (vp *VipsProcessor) Init() error {
	// noop for this processor
	return nil
}

func (vp *VipsProcessor) NewImageHandler(s string) (Handler, error) {
	return &VipsHandler{}, nil
}
