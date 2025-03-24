package background

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/qbarrand/quba.fr/internal/metadata"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	baseURL      string
	logger       logrus.FieldLogger
	mdb          *metadata.DB
	validFormats map[string]bool
}

func NewHandler(imgDir string, baseURL string, logger logrus.FieldLogger) (*Handler, error) {
	dbPath := filepath.Join(imgDir, "metadata.db")

	mdb, err := metadata.OpenDB(dbPath, true)
	if err != nil {
		return nil, fmt.Errorf("could not open the database: %v", err)
	}

	h := Handler{
		baseURL: baseURL,
		logger:  logger,
		mdb:     mdb,
		validFormats: map[string]bool{
			"jpg":  true,
			"webp": true,
		},
	}

	return &h, nil
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		err           error
		format        string
		name          *string
		width, height *int
	)

	logger := h.logger

	if err = r.ParseForm(); err != nil {
		logger.WithError(err).Error("could not parse the query string")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if nameStr := r.Form.Get("name"); nameStr != "" {
		nameStr = removeLineBreaks(nameStr)
		name = &nameStr

		logger = logger.WithField("name", name)
	}

	if widthStr := r.Form.Get("width"); widthStr != "" {
		widthInt, err := strconv.Atoi(widthStr)
		if err != nil {
			logger.WithError(err).Error("could not parse width")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		width = &widthInt
		logger = logger.WithField("width", *width)
	}

	if heightStr := r.Form.Get("height"); heightStr != "" {
		heightStr = removeLineBreaks(heightStr)

		heightInt, err := strconv.Atoi(heightStr)
		if err != nil {
			logger.WithError(err).Error("could not parse height")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		height = &heightInt
		logger = logger.WithField("height", *height)
	}

	if format = r.Form.Get("format"); !h.validFormats[format] {
		logger.
			WithField(
				"format",
				removeLineBreaks(format),
			).Error("invalid format")

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	logger = logger.WithField("format", format)

	ir, err := h.mdb.GetImage(r.Context(), name, width, height, format)
	if err != nil {
		logger.WithError(err).Error("no suitable image found")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := struct {
		URL       string `json:"url"`
		Name      string `json:"name"`
		Date      string `json:"date"`
		Location  string `json:"location"`
		MainColor string `json:"mainColor"`
	}{
		URL:       filepath.Join(h.baseURL, ir.Filename),
		Name:      ir.Name,
		Date:      ir.Date,
		Location:  ir.Location,
		MainColor: ir.MainColor,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&body); err != nil {
		logger.WithError(err).Error("Could not encode JSON")
	}
}

func removeLineBreaks(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")

	return s
}
