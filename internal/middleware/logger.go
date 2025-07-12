package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// Middleware returns an http.Handler that logs each request.
func Logger(next http.Handler) http.Handler {
	l := slog.Default()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: 200}

		start := time.Now()
		next.ServeHTTP(sw, r)
		elapsed := time.Since(start)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		l.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", sw.status,
			"duration", elapsed,
			"remote", ip,
		)
	})
}
