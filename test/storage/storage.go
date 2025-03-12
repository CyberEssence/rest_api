// storage/storage.go
package storage

import (
	"fmt"
	"sync"
	"test/models"
)

type InMemoryStorage struct {
	tasks  map[int]*models.Task
	lastID int
	mu     sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[int]*models.Task),
	}
}

func (s *InMemoryStorage) CreateTask(title, description string) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	task := &models.Task{
		ID:          s.lastID,
		Title:       title,
		Description: description,
		Completed:   false,
	}
	s.tasks[s.lastID] = task
	return task, nil
}

func (s *InMemoryStorage) GetAllTasks() ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *InMemoryStorage) GetTask(id int) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("задача с ID %d не найдена", id)
	}
	return task, nil
}

func (s *InMemoryStorage) UpdateTask(id int, title, description string, completed bool) (*models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("задача с ID %d не найдена", id)
	}

	task.Title = title
	task.Description = description
	task.Completed = completed
	return task, nil
}

func (s *InMemoryStorage) DeleteTask(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}
	delete(s.tasks, id)
	return nil
}
