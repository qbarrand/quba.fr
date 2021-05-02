package imagefs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/qbarrand/quba.fr/data/images"
	"github.com/sirupsen/logrus"
)

func ImageLister(mfs images.MetadataFS, logger logrus.FieldLogger) (http.HandlerFunc, error) {
	var buf bytes.Buffer

	entries, err := fs.ReadDir(mfs, ".")
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
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(b); err != nil {
			logger.WithError(err).Error("Could not write the list of images")
		}
	}

	return handler, nil
}
