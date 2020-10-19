package image

import (
	"errors"
)

type Format uint

const (
	JPEG Format = iota
	Webp
)

var ErrNotAcceptable = errors.New("no acceptable MIME type found")

func AcceptHeaderToFormat(accept []string) (Format, string, error) {
	for _, mimeType := range accept {
		switch mimeType {
		case "image/jpeg":
			return JPEG, mimeType, nil
		case "image/webp":
			return Webp, mimeType, nil
		}
	}

	return 0, "", ErrNotAcceptable
}
