//go:generate go run github.com/golang/mock/mockgen -source processor.go -destination mock_image/processor.go Handler,Processor

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
}