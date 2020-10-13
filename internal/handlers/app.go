package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type AppOptions struct {
	AssetsDir string
	LastMod   time.Time
}

func NewApp(opts *AppOptions, logger logrus.FieldLogger) (*mux.Router, error) {
	lastMod := opts.LastMod.Format("2006-01-02")

	sitemap, err := newSitemap(lastMod)
	if err != nil {
		return nil, fmt.Errorf("could not create the sitemap handler: %v", err)
	}

	router := mux.NewRouter()

	subRouter := router.Methods(http.MethodGet).Subrouter()
	subRouter.Handle("/sitemap.xml", sitemap)
	subRouter.Handle("/", http.FileServer(http.Dir("webroot")))

	return router, nil
}
