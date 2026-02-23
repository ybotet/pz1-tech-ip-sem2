package service

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"-"`
}

type TaskService struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[string]Task),
	}
}

func (s *TaskService) Create(task Task) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	task.ID = fmt.Sprintf("t_%03d", len(s.tasks)+1)
	task.CreatedAt = time.Now()
	s.tasks[task.ID] = task
	return task
}

func (s *TaskService) GetAll() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *TaskService) GetByID(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	task, ok := s.tasks[id]
	return task, ok
}

func (s *TaskService) Update(id string, updated Task) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task, ok := s.tasks[id]
	if !ok {
		return Task{}, false
	}
	if updated.Title != "" {
		task.Title = updated.Title
	}
	if updated.Description != "" {
		task.Description = updated.Description
	}
	if updated.DueDate != "" {
		task.DueDate = updated.DueDate
	}
	task.Done = updated.Done
	s.tasks[id] = task
	return task, true
}

func (s *TaskService) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; ok {
		delete(s.tasks, id)
		return true
	}
	return false
}
