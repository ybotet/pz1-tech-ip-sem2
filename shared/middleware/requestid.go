package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"
const RequestIDHeader = "X-Request-ID"

// RequestIDMiddleware genera o propaga un X-Request-ID y lo añade al contexto
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		// Añadir al contexto
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		// Añadir al header de respuesta para trazabilidad
		w.Header().Set(RequestIDHeader, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
