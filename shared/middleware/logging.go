package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware registra método, ruta, duración y ID de solicitud
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := r.Header.Get(RequestIDHeader)
		log.Printf("[%s] %s %s %s", requestID, r.Method, r.URL.Path, time.Since(start))
		next.ServeHTTP(w, r)
	})
}
