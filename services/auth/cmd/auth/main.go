package main

import (
	"auth/internal/handlers"
	"auth/internal/service"
	"log"
	"net/http"
	"os"

	// "auth/shared/middleware"
	"github.com/ybotet/pz1-tech-ip-sem2/shared/middleware"
)

func main() {
	// Configuración desde variables de entorno
	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8081"
	}

	// Inicializar dependencias
	authService := service.NewAuthService()
	authHandlers := handlers.NewAuthHandlers(authService)

	// Configurar rutas
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/auth/login", authHandlers.Login)
	mux.HandleFunc("GET /v1/auth/verify", authHandlers.Verify)

	// Aplicar middlewares globales
	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.LoggingMiddleware(handler)

	log.Printf("Auth service iniciado en puerto %s", port)
	log.Printf("Endpoints disponibles:")
	log.Printf("  POST http://localhost:%s/v1/auth/login", port)
	log.Printf("  GET http://localhost:%s/v1/auth/verify", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
