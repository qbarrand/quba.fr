package handlers

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	boolCache struct {
		lastCheck time.Time
		valid     bool

		m sync.Mutex
	}

	healthz struct {
		cache       *boolCache
		dnsQueryier func(string) ([]string, error)
		logger      logrus.FieldLogger
	}
)

func newHealthz(logger logrus.FieldLogger) *healthz {
	return &healthz{
		cache:       &boolCache{},
		dnsQueryier: net.LookupTXT,
		logger:      logger,
	}
}

func (h *healthz) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const (
		fqdn     = "ping.quba.fr"
		interval = 120 * time.Second
	)

	h.logger.
		WithField("time", h.cache.lastCheck).
		Debug("Checking duration since last lookup")

	now := time.Now()
	elapsed := time.Now().Sub(h.cache.lastCheck)

	if elapsed < interval {
		h.logger.WithField("elapsed", elapsed).Debug("Using cache")
	} else {
		h.logger.WithField("elapsed", elapsed).Debug("Cache too old")

		records, err := h.dnsQueryier(fqdn)
		if err != nil {
			h.logger.WithError(err).Error("Error while running the query")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(records) != 1 {
			h.logger.WithField("query", fqdn+"/TXT").Error("Not enough records")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		const expected = "quentin@quba.fr"
		got := records[0]

		h.cache.m.Lock()
		h.cache.valid = got == expected

		if !h.cache.valid {
			h.logger.
				WithFields(logrus.Fields{
					"expected": expected,
					"got":      got,
				}).
				Error("Unexpected record value")
		}

		h.cache.lastCheck = now
		h.cache.m.Unlock()
	}

	if !h.cache.valid {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
