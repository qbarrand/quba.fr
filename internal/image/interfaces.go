package image

import (
	"context"
	"io"
)

type Handler interface {
	io.WriterTo

	Resize(context.Context, int, int) error
}

type Processor interface {
	Init() error
	NewImageHandler(string) (*Handler, error)
}
