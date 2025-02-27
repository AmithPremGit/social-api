package middleware

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"
)

// LoggingMiddleware is a custom logging middleware
func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Start timer
			start := time.Now()

			// Create a response wrapper to capture status code
			lrw := newLoggingResponseWriter(w)

			// Process the request
			next.ServeHTTP(lrw, r)

			// Calculate request processing time
			duration := time.Since(start)

			// Log the request details
			logger.Printf(
				"REQUEST: method=%s path=%s status=%d duration=%s ip=%s",
				r.Method,
				r.URL.Path,
				lrw.statusCode,
				duration,
				r.RemoteAddr,
			)
		})
	}
}

// loggingResponseWriter is a wrapper for http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newLoggingResponseWriter creates a new logging response writer
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default to 200 OK
	}
}

// WriteHeader captures the status code and passes it to the wrapped ResponseWriter
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Implement the http.Pusher interface if the wrapped ResponseWriter implements it
func (lrw *loggingResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := lrw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Implement the http.Flusher interface if the wrapped ResponseWriter implements it
func (lrw *loggingResponseWriter) Flush() {
	if flusher, ok := lrw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Implement the http.Hijacker interface if the wrapped ResponseWriter implements it
func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := lrw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
