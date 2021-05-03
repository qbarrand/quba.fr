package image

import (
	"context"
)

type Handler interface {
	Bytes() ([]byte, error)
	Destroy() error
	Resize(context.Context, int, int) error
	SetFormat(Format) error
}

type Processor interface {
	BestFormats() []Format
	Destroy() error
	Init() error
	NewImageHandler(string) (Handler, error)
	HandlerFromBytes([]byte) (Handler, error)
}
