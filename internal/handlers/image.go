package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	image3 "github.com/qbarrand/quba.fr/internal/image"
	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/pkg/httputils"
)

type image struct {
	logger    logrus.FieldLogger
	metaDB    image3.MetaDB
	path      string
	processor image3.Processor
}

func newImage(processor image3.Processor, path string, metaDB image3.MetaDB, logger logrus.FieldLogger) (*image, error) {
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

	// Get the accepted MIME types
	acceptedTypes, err := getMIMETypes(r.Header.Values("Accept"))
	if err != nil {
		logger.WithError(err).Error("Could not parse the Accept headers")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	format, err := getBestFormat(i.processor.BestFormats(), acceptedTypes)
	if err != nil {
		logger.WithError(err).Error("Could not find the best MIME type")
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if err := handler.SetFormat(format); err != nil {
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
	headers.Set("Content-Type", string(format))
	headers.Set("ETag", hex.EncodeToString(h.Sum(nil)))
	headers.Set("X-Quba-Date", strconv.FormatInt(meta.Date.Unix(), 10))
	headers.Set("X-Quba-Location", meta.Location)
	headers.Set("X-Quba-MainColor", meta.MainColor)

	n, err := w.Write(buf)
	if err != nil {
		// do not call fail() as we've already written headers
		logger.WithError(err).Error("Could not write the resulting image")
		return
	}

	logger.WithField("bytes", n).Debug("Finished writing image")
}

func newImageLister(metaDB image3.MetaDB, logger logrus.FieldLogger) (http.HandlerFunc, error) {
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

func getBestFormat(serverFormats []image3.Format, clientTypes []string) (image3.Format, error) {
	clientFmtMap := make(map[image3.Format]bool, len(clientTypes))

	for _, t := range clientTypes {
		f, err := image3.FormatFromMIMEType(t)
		if err != nil {
			// This MIME type is not a recognized format
			continue
		}

		clientFmtMap[f] = true
	}

	for _, f := range serverFormats {
		if clientFmtMap[f] {
			return f, nil
		}
	}

	return "", errors.New("no acceptable format found")
}

// getMIMETypes parses a slice of Accept headers and returns a slice of all the types.
// It returns an error if a header element could not be parsed.
func getMIMETypes(acceptHeaders []string) ([]string, error) {
	types := make([]string, 0)

	for _, header := range acceptHeaders {
		for _, i := range strings.Split(header, ",") {
			mt, _, err := mime.ParseMediaType(i)
			if err != nil {
				return nil, fmt.Errorf("could not parse header item %s: %v", i, err)
			}

			types = append(types, mt)
		}
	}

	return types, nil
}
