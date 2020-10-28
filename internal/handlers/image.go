package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"

	img "github.com/qbarrand/quba.fr/internal/image"
	"github.com/qbarrand/quba.fr/pkg/httputils"
)

type image struct {
	logger    logrus.FieldLogger
	metaDB    img.MetaDB
	path      string
	processor img.Processor
}

func newImage(processor img.Processor, path string, metaDB img.MetaDB, logger logrus.FieldLogger) (*image, error) {
	if err := processor.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize the processor: %w", err)
	}

	i := image{
		logger:    logger,
		metaDB:    metaDB,
		path:      path,
		processor: processor,
	}

	return &i, nil
}

func (i *image) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := filepath.Base(r.URL.Path)

	path := filepath.Join(i.path, name)

	logger := i.logger.WithFields(logrus.Fields{
		"path":       path,
		"request-id": httputils.GetRequestID(r),
	})

	logger.Debug("Serving image")

	fail := func(msg string, wrapped error) {
		logger.WithError(wrapped).Error(msg)
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

	h := fnv.New32()
	h.Write(buf) // per the docs: never returns an error

	meta, err := i.metaDB.GetMetadata(name)
	if err != nil {
		fail("Could not get the image metadata", err)
		return
	}

	headers := w.Header()

	headers.Set("Content-Length", strconv.Itoa(len(buf)))
	headers.Set("Content-Type", mimeType)
	headers.Set("ETag", hex.EncodeToString(h.Sum(nil)))
	headers.Set("X-Quba-Date", strconv.FormatInt(meta.Date.Unix(), 10))
	headers.Set("X-Quba-Location", meta.Location)

	n, err := w.Write(buf)
	if err != nil {
		// do not call fail() as we've already written headers
		logger.WithError(err).Error("Could not write the resulting image")
		return
	}

	logger.WithField("bytes", n).Debug("Finished writing image")
}

func newImageLister(metaDB img.MetaDB, logger logrus.FieldLogger) (http.HandlerFunc, error) {
	var buf bytes.Buffer

	allNames, err := metaDB.AllNames()
	if err != nil {
		return nil, fmt.Errorf("could not get a list of image names: %w", err)
	}

	if err := json.NewEncoder(&buf).Encode(allNames); err != nil {
		return nil, fmt.Errorf("could not generate the corresponding JSON: %v", err)
	}

	b := buf.Bytes()

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write(b); err != nil {
			logger.WithError(err).Error("Could not write the list of images")
		}
	}

	return handlerFunc, nil
}
