// Package imgpro offers Image Processing interfaces as well as a few common implementations.
package imgpro

import (
	"context"
	"fmt"
)

type Handler interface {
	Bytes() ([]byte, error)
	Destroy() error
	Resize(context.Context, int, int) error
	SetFormat(Format) error
	StripMetadata() error
}

type Processor interface {
	fmt.Stringer

	BestFormats() []Format
	Destroy() error
	Init(int) error
	HandlerFromBytes([]byte) (Handler, error)
}
