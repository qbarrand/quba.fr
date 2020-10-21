package httputils

import "net/http"

type StatusCodeResponseWriter struct {
	statusCode int
	wrapped    http.ResponseWriter
}

func NewStatusCodeResponseWriter(wrapped http.ResponseWriter) *StatusCodeResponseWriter {
	return &StatusCodeResponseWriter{wrapped: wrapped}
}

func (s *StatusCodeResponseWriter) Header() http.Header {
	return s.wrapped.Header()
}

func (s *StatusCodeResponseWriter) StatusCode() int {
	return s.statusCode
}

func (s *StatusCodeResponseWriter) Write(b []byte) (int, error) {
	// Respect the http.ResponseWriter interface's behaviour as documented
	if s.statusCode == 0 {
		s.WriteHeader(http.StatusOK)
	}

	return s.wrapped.Write(b)
}

func (s *StatusCodeResponseWriter) WriteHeader(statusCode int) {
	s.statusCode = statusCode
	s.wrapped.WriteHeader(statusCode)
}
