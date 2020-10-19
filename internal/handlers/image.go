package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"

	img "github.com/qbarrand/quba.fr/internal/image"
)

type image struct {
	logger    logrus.FieldLogger
	path      string
	processor img.Processor
}

func newImage(processor img.Processor, path string, logger logrus.FieldLogger) (*image, error) {
	if err := processor.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize the processor: %w", err)
	}

	i := image{
		logger:    logger,
		path:      path,
		processor: processor,
	}

	return &i, nil
}

func (i *image) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(
		i.path,
		filepath.Base(r.URL.Path),
	)

	i.logger.WithField("path", path).Debug("Serving image")

	fail := func(msg string, wrapped error) {
		i.logger.WithError(wrapped).Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
	}

	handler, err := i.processor.NewImageHandler(path)
	if err != nil {
		fail("Could not create an image handler", err)
		return
	}
	defer handler.Destroy()

	// Get the format
	f, mimeType, err := img.AcceptHeaderToFormat(r.Header.Values("Accept"))
	if err != nil {
		if errors.Is(img.ErrNotAcceptable, err) {
			// TODO find something clever in case the image is not JPEG
			i.logger.WithError(err).Debug("No acceptable format: using JPEG")
			f = img.JPEG
		} else {
			fail("Error while determining the accepted format", err)
			return
		}
	}

	w.Header().Set("Content-Type", mimeType)

	if err := handler.SetFormat(f); err != nil {
		fail("Could not set the format", err)
		return
	}

	widthStr := r.FormValue("width")

	if widthStr != "" {
		width, err := strconv.Atoi(widthStr)
		if err != nil {
			fail("Could not convert the width parameter to integer", err)
			return
		}

		if err := handler.Resize(r.Context(), width, 0); err != nil {
			fail("Could not resize the image", err)
			return
		}
	}

	buf, err := handler.Bytes()
	if err != nil {
		fail("Could not get the image bytes", err)
		return
	}

	n, err := bytes.NewReader(buf).WriteTo(w)
	if err != nil {
		fail("Could not write the resulting image", err)
		return
	}

	i.logger.WithField("bytes", n).Debug("Finished writing image")
}
