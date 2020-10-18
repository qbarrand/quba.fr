package image

import (
	"context"
	"io"
)

type Handler interface {
	io.WriterTo

	Resize(context.Context, int, int) error
	SetFormat(Format) error
}

type Processor interface {
	Destroy() error
	Init() error
	NewImageHandler(string) (Handler, error)
}
