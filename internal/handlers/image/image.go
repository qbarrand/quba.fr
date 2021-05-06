package image

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/qbarrand/quba.fr/data/images"
	"github.com/qbarrand/quba.fr/internal/imgpro"
	"github.com/qbarrand/quba.fr/pkg/httputils"
)

type Image struct {
	mfs       images.MetadataFS
	logger    logrus.FieldLogger
	processor imgpro.Processor
}

func New(processor imgpro.Processor, mfs images.MetadataFS, logger logrus.FieldLogger) (*Image, error) {
	if err := processor.Init(); err != nil {
		return nil, fmt.Errorf("could not initialize the processor: %w", err)
	}

	i := Image{
		mfs:       mfs,
		logger:    logger,
		processor: processor,
	}

	return &i, nil
}

func (i *Image) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/")

	logger := i.logger.WithFields(logrus.Fields{
		"name":       name,
		"request-id": httputils.GetRequestID(r),
	})

	logger.Debug("Serving image")

	fail := func(msg string, wrapped error) {
		logger.WithError(wrapped).Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fd, meta, err := i.mfs.OpenWithMetadata(name)
	if err != nil {
		fail("Could not open the image", err)
		return
	}
	defer fd.Close()

	b, err := io.ReadAll(fd)
	if err != nil {
		fail("Could not read the image", err)
		return
	}

	handler, err := i.processor.HandlerFromBytes(b)
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

	width, height, err := getDimensions(r)
	if err != nil {
		logger.WithError(err).Error("Invalid width or height")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if width != 0 || height != 0 {
		if err := handler.Resize(r.Context(), width, height); err != nil {
			fail("Could not resize the image", err)
			return
		}
	}

	if err := handler.StripMetadata(); err != nil {
		fail("Could not strip the image's metadata", err)
		return
	}

	buf, err := handler.Bytes()
	if err != nil {
		fail("Could not get the image bytes", err)
		return
	}

	h := fnv.New32()
	h.Write(buf) // per the docs: never returns an error

	headers := w.Header()

	headers.Set("Content-Length", strconv.Itoa(len(buf)))
	headers.Set("Content-Type", string(format))
	headers.Set("ETag", hex.EncodeToString(h.Sum(nil)))
	headers.Set("X-Quba-Date", strconv.FormatInt(meta.Date.Unix(), 10))
	headers.Set("X-Quba-Location", meta.Location)
	headers.Set("X-Quba-Main-Color", meta.MainColor)

	rs := bytes.NewReader(buf)

	modTime := time.Time{}

	stat, err := fd.Stat()
	if err != nil {
		logger.WithError(err).Error("Could not stat(); using zero time")
	} else {
		modTime = stat.ModTime()
	}

	http.ServeContent(w, r, "", modTime, rs)
}

func getBestFormat(serverFormats []imgpro.Format, clientTypes []string) (imgpro.Format, error) {
	clientFmtMap := make(map[imgpro.Format]bool, len(clientTypes))

	for _, t := range clientTypes {
		f, err := imgpro.FormatFromMIMEType(t)
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

func getDimensions(req *http.Request) (int, int, error) {
	width := 0
	height := 0

	var err error

	if widthStr := req.FormValue("width"); widthStr != "" {
		width, err = strconv.Atoi(widthStr)
		if err != nil {
			return width, height, fmt.Errorf("invalid width: %v", err)
		}
	}

	if heightStr := req.FormValue("height"); heightStr != "" {
		height, err = strconv.Atoi(heightStr)
		if err != nil {
			return width, height, fmt.Errorf("invalid height: %v", err)
		}
	}

	return width, height, nil
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
