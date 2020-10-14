package handlers

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	img "github.com/qbarrand/quba.fr/internal/image"
)

const requestIDKey = "request-id"

func getRequestID(r *http.Request) string {
	return r.Context().Value(requestIDKey).(string)
}

type AppOptions struct {
	ImageProcessor img.Processor
	LastMod        time.Time
	WebRootDir     string
}

func loggingMiddleware(logger logrus.FieldLogger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id string

		// uuid.NewRandom does not always work, especially if the random generator
		// runs out of randomness.
		u, err := uuid.NewRandom()
		if err != nil {
			id = fmt.Sprintf("<error: %v>", err)
		} else {
			id = u.String()
		}

		logger.
			WithFields(logrus.Fields{
				"id":     id,
				"method": r.Method,
				"remote": r.RemoteAddr,
				"url":    r.URL.String(),
			}).
			Info("New request")

		ctx := context.WithValue(r.Context(), requestIDKey, id)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	}
}

func NewApp(opts *AppOptions, logger logrus.FieldLogger) (http.HandlerFunc, error) {
	const handlerKey = "handler"

	sitemap, err := newSitemap(opts.LastMod, logger.WithField(handlerKey, "sitemap"))
	if err != nil {
		return nil, fmt.Errorf("could not create the sitemap handler: %v", err)
	}

	const subdir = "images"

	imagesPath := filepath.Join(opts.WebRootDir, subdir)

	image, err := newImage(opts.ImageProcessor, imagesPath, logger.WithField(handlerKey, "image"))
	if err != nil {
		return nil, fmt.Errorf("could not initialize the image handler: %v", err)
	}

	router := mux.NewRouter()

	subRouter := router.Methods(http.MethodGet).Subrouter()
	subRouter.Handle("/sitemap.xml", sitemap)
	subRouter.PathPrefix("/images").Handler(image)
	subRouter.PathPrefix("/").Handler(
		http.FileServer(
			http.Dir(opts.WebRootDir),
		),
	)

	// Intercept all requests, then forward them to the router.
	// We use this instead of a middleware, as those are only hit when the router
	// has a match (not for 404, 405 etc).
	return loggingMiddleware(logger, router), nil
}
