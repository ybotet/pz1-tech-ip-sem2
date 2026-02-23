package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"tasks/internal/client/authclient"
	"tasks/internal/service"
	"time"

	"github.com/ybotet/pz1-tech-ip-sem2/shared/middleware"
)

type TaskHandlers struct {
	taskService *service.TaskService
	authClient  *authclient.AuthClient
}

func NewTaskHandlers(ts *service.TaskService, ac *authclient.AuthClient) *TaskHandlers {
	return &TaskHandlers{taskService: ts, authClient: ac}
}

// verifyTokenMiddleware es un helper para verificar el token antes de cada operación
func (h *TaskHandlers) verifyTokenFromRequest(r *http.Request) (bool, string, int, interface{}) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return false, "", http.StatusUnauthorized, map[string]string{"error": "missing authorization header"}
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return false, "", http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"}
	}
	token := parts[1]

	requestID := r.Context().Value(middleware.RequestIDKey).(string)

	// Crear un contexto con timeout para la llamada a Auth
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	valid, subject, err := h.authClient.VerifyToken(ctx, token, requestID)
	if err != nil {
		// Error de comunicación con Auth (timeout, caída, etc.) -> Fail Closed (500)
		return false, "", http.StatusInternalServerError, map[string]string{"error": "authorization service unavailable"}
	}
	if !valid {
		return false, "", http.StatusUnauthorized, map[string]string{"error": "invalid token"}
	}

	// Token válido, podemos usar el subject si es necesario
	_ = subject
	return true, "", http.StatusOK, nil
}

// AuthMiddleware para envolver handlers protegidos
func (h *TaskHandlers) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valid, _, status, errBody := h.verifyTokenFromRequest(r)
		if !valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(errBody)
			return
		}
		next(w, r)
	}
}

// CreateTask maneja POST /v1/tasks
func (h *TaskHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task service.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	created := h.taskService.Create(task)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// GetTasks maneja GET /v1/tasks
func (h *TaskHandlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.taskService.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTask maneja GET /v1/tasks/{id}
func (h *TaskHandlers) GetTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	task, ok := h.taskService.GetByID(id)
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// UpdateTask maneja PATCH /v1/tasks/{id}
func (h *TaskHandlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	var updated service.Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	task, ok := h.taskService.Update(id, updated)
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask maneja DELETE /v1/tasks/{id}
func (h *TaskHandlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	ok := h.taskService.Delete(id)
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
