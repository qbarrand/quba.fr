package image

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/sirupsen/logrus"
)

// Lister returns a http.HandlerFunc that lists all img-src available in mfs.
func Lister(fsys fs.FS, logger logrus.FieldLogger) (http.HandlerFunc, error) {
	var buf bytes.Buffer

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("could not list files: %v", err)
	}

	names := make([]string, 0, len(entries))

	for _, de := range entries {
		names = append(names, de.Name())
	}

	if err := json.NewEncoder(&buf).Encode(names); err != nil {
		return nil, fmt.Errorf("could not serialize the slice of names into JSON: %v", err)
	}

	b := buf.Bytes()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			logger.WithError(err).Error("Could not write the list of img-src")
		}
	}

	return handler, nil
}
