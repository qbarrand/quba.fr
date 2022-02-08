// +build magick

package main

import "github.com/qbarrand/quba.fr/internal/imgpro"

func init() {
	processorMap["magick"] = &imgpro.ImageMagickProcessor{}
}
