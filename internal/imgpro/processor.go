// Package imgpro offers Image Processing interfaces as well as a few common implementations.
package imgpro

import (
	"context"
)

type Handler interface {
	Bytes() ([]byte, error)
	Destroy() error
	Resize(context.Context, int, int) error
	SetFormat(Format) error
	StripMetadata() error
}

type Processor interface {
	BestFormats() []Format
	Destroy() error
	Init() error
	NewImageHandler(string) (Handler, error)
	HandlerFromBytes([]byte) (Handler, error)
}
