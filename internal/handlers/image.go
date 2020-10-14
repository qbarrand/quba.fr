package handlers

import (
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

	http.ServeFile(w, r, path)
}
