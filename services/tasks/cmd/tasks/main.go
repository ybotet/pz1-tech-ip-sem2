package main

import (
	"log"
	"net/http"
	"os"
	"tasks/internal/client/authclient"
	"tasks/internal/handlers"
	"tasks/internal/service"

	"github.com/ybotet/pz1-tech-ip-sem2/shared/middleware"
)

func main() {
	tasksPort := os.Getenv("TASKS_PORT")
	if tasksPort == "" {
		tasksPort = "8082"
	}
	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
	}

	taskService := service.NewTaskService()
	authClient := authclient.NewAuthClient(authBaseURL)
	taskHandlers := handlers.NewTaskHandlers(taskService, authClient)

	mux := http.NewServeMux()
	// Rutas protegidas por el middleware de autenticación
	mux.HandleFunc("POST /v1/tasks", taskHandlers.AuthMiddleware(taskHandlers.CreateTask))
	mux.HandleFunc("GET /v1/tasks", taskHandlers.AuthMiddleware(taskHandlers.GetTasks))
	mux.HandleFunc("GET /v1/tasks/", taskHandlers.AuthMiddleware(taskHandlers.GetTask))
	mux.HandleFunc("PATCH /v1/tasks/", taskHandlers.AuthMiddleware(taskHandlers.UpdateTask))
	mux.HandleFunc("DELETE /v1/tasks/", taskHandlers.AuthMiddleware(taskHandlers.DeleteTask))

	// Middlewares globales
	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.LoggingMiddleware(handler)

	log.Printf("Tasks service iniciado en puerto %s", tasksPort)
	if err := http.ListenAndServe(":"+tasksPort, handler); err != nil {
		log.Fatal(err)
	}
}
