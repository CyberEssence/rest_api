// Package storage предоставляет реализацию хранилища задач в памяти
package storage

import (
	"fmt"
	"sync"
	"test/models"
)

// InMemoryStorage реализует хранилище задач в памяти с поддержкой конкурентного доступа
type InMemoryStorage struct {
	tasks  map[int]*models.Task // Хранилище задач
	lastID int                  // Последний использованный ID
	mu     sync.RWMutex         // Мьютекс для синхронизации доступа
}

// NewInMemoryStorage создает новое хранилище задач в памяти
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		tasks: make(map[int]*models.Task),
	}
}

// CreateTask создает новую задачу в хранилище
//
// Args:
//
//	title: название задачи
//	description: описание задачи
//
// Returns:
//
//	*models.Task: созданная задача
//	error: ошибка при создании задачи
func (s *InMemoryStorage) CreateTask(title, description string) (*models.Task, error) {
	// Блокировка на запись для атомарного создания задачи
	s.mu.Lock()
	defer s.mu.Unlock()

	// Генерация нового ID
	s.lastID++

	// Создание новой задачи
	task := &models.Task{
		ID:          s.lastID,
		Title:       title,
		Description: description,
		Completed:   false,
	}

	// Сохранение задачи в хранилище
	s.tasks[s.lastID] = task
	return task, nil
}

// GetAllTasks возвращает список всех задач из хранилища
//
// Returns:
//
//	[]*models.Task: список всех задач
//	error: ошибка при получении задач
func (s *InMemoryStorage) GetAllTasks() ([]*models.Task, error) {
	// Блокировка на чтение для безопасного получения всех задач
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Создание нового среза для хранения задач
	tasks := make([]*models.Task, 0, len(s.tasks))

	// Копирование всех задач в новый срез
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetTask возвращает задачу по ID
//
// Args:
//
//	id: ID задачи
//
// Returns:
//
//	*models.Task: найденная задача
//	error: ошибка при поиске задачи
func (s *InMemoryStorage) GetTask(id int) (*models.Task, error) {
	// Блокировка на чтение для безопасного получения задачи
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Поиск задачи по ID
	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("задача с ID %d не найдена", id)
	}

	return task, nil
}

// UpdateTask обновляет существующую задачу
//
// Args:
//
//	id: ID задачи
//	title: новое название задачи
//	description: новое описание задачи
//	completed: новый статус выполнения
//
// Returns:
//
//	*models.Task: обновленная задача
//	error: ошибка при обновлении задачи
func (s *InMemoryStorage) UpdateTask(id int, title, description string, completed bool) (*models.Task, error) {
	// Блокировка на запись для атомарного обновления задачи
	s.mu.Lock()
	defer s.mu.Unlock()

	// Поиск задачи по ID
	task, exists := s.tasks[id]
	if !exists {
		return nil, fmt.Errorf("задача с ID %d не найдена", id)
	}

	// Обновление полей задачи
	task.Title = title
	task.Description = description
	task.Completed = completed

	return task, nil
}

// DeleteTask удаляет задачу из хранилища
//
// Args:
//
//	id: ID задачи для удаления
//
// Returns:
//
//	error: ошибка при удалении задачи
func (s *InMemoryStorage) DeleteTask(id int) error {
	// Блокировка на запись для атомарного удаления задачи
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверка существования задачи
	if _, exists := s.tasks[id]; !exists {
		return fmt.Errorf("задача с ID %d не найдена", id)
	}

	// Удаление задачи из хранилища
	delete(s.tasks, id)
	return nil
}
