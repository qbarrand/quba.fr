package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

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

	// Get the format
	f, err := img.AcceptHeaderToFormat(r.Header.Values("Accept"))
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

	_ = f
	_ = handler
}
