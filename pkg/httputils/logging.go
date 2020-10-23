package httputils

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const requestIDKey = "request-id"

func GetRequestID(r *http.Request) string {
	v := r.Context().Value(requestIDKey)

	if v == nil {
		return "<nil>"
	}

	if id, ok := v.(string); ok {
		return id
	}

	if id, ok := v.(int); ok {
		return strconv.Itoa(id)
	}

	return "<error>"
}

func LoggingMiddleware(logger logrus.FieldLogger, next http.Handler) http.HandlerFunc {
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

		logger = logger.WithField("id", id)

		logger.
			WithFields(logrus.Fields{
				"method": r.Method,
				"remote": r.RemoteAddr,
				"url":    r.URL.String(),
			}).
			Info("New request")

		ctx := context.WithValue(r.Context(), requestIDKey, id)

		scrw := NewStatusCodeResponseWriter(w)

		next.ServeHTTP(
			scrw,
			r.WithContext(ctx),
		)

		logger.
			WithField("status", scrw.StatusCode()).
			Info("Finished serving request")
	}
}
