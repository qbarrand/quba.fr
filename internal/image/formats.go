package image

import (
	"errors"
)

type Format string

const (
	JPEG Format = "image/jpeg"
	Webp Format = "image/webp"
)

var (
	ErrNotAcceptable = errors.New("no acceptable MIME type found")

	mimeToFormat = map[string]Format{
		string(JPEG): JPEG,
		string(Webp): Webp,
	}
)

func FormatFromMIMEType(mimeType string) (Format, error) {
	f := mimeToFormat[mimeType]

	if f == "" {
		return f, ErrNotAcceptable
	}

	return f, nil
}
