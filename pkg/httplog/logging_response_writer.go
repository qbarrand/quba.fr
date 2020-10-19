package httplog

import "net/http"

type LoggingResponseWriter struct {
	statusCode int
	wrapped    http.ResponseWriter
}

func NewLoggingResponseWriter(wrapped http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{wrapped: wrapped}
}

func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.wrapped.Header()
}

func (lrw *LoggingResponseWriter) StatusCode() int {
	return lrw.statusCode
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	// Respect the http.ResponseWriter interface's behaviour as documented
	if lrw.statusCode == 0 {
		lrw.WriteHeader(http.StatusOK)
	}

	return lrw.wrapped.Write(b)
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.wrapped.WriteHeader(statusCode)
}
