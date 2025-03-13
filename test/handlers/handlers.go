// Package handlers предоставляет HTTP обработчики для работы с задачами
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test/storage"
)

// SetupHandlers настраивает маршрутизатор HTTP с обработчиками для работы с задачами
func SetupHandlers(storage *storage.InMemoryStorage) *http.ServeMux {
	mux := http.NewServeMux()

	// Регистрация обработчиков для /tasks
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateTaskHandler(w, r, storage)
		case http.MethodGet:
			GetAllTasksHandler(w, r, storage)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// Регистрация обработчиков для /tasks/{id}
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/tasks/" {
			http.Error(w, "ID не указан", http.StatusBadRequest)
			return
		}

		idStr := r.URL.Path[len("/tasks/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный формат ID", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			GetTaskHandler(w, r, storage, id)
		case http.MethodPut:
			UpdateTaskHandler(w, r, storage, id)
		case http.MethodDelete:
			DeleteTaskHandler(w, r, storage, id)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	return mux
}

// CreateTaskHandler создает новую задачу
// POST /tasks
//
// Запрос:
//
//	{
//	  "title": "Название задачи",
//	  "description": "Описание задачи"
//	}
//
// Ответ:
//
//	{
//	  "id": 1,
//	  "title": "Название задачи",
//	  "description": "Описание задачи",
//	  "completed": false
//	}
func CreateTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage) {
	var taskData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	// Декодирование JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Валидация входных данных
	if taskData.Title == "" || taskData.Description == "" {
		http.Error(w, "Title и Description обязательны", http.StatusBadRequest)
		return
	}

	// Создание задачи в хранилище
	task, err := storage.CreateTask(taskData.Title, taskData.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возврат созданной задачи с кодом 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetAllTasksHandler возвращает список всех задач
// GET /tasks
//
// Ответ:
// [
//
//	{
//	  "id": 1,
//	  "title": "Задача 1",
//	  "description": "Описание 1",
//	  "completed": false
//	}
//
// ]
func GetAllTasksHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage) {
	tasks, err := storage.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskHandler возвращает задачу по ID
// GET /tasks/{id}
//
// Ответ:
//
//	{
//	  "id": 1,
//	  "title": "Задача 1",
//	  "description": "Описание 1",
//	  "completed": false
//	}
func GetTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage, id int) {
	task, err := storage.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

// UpdateTaskHandler обновляет существующую задачу
// PUT /tasks/{id}
//
// Запрос:
//
//	{
//	  "title": "Новое название",
//	  "description": "Новое описание",
//	  "completed": true
//	}
//
// Ответ:
//
//	{
//	  "id": 1,
//	  "title": "Новое название",
//	  "description": "Новое описание",
//	  "completed": true
//	}
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage, id int) {
	var taskData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	// Декодирование JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновление задачи в хранилище
	task, err := storage.UpdateTask(id, taskData.Title, taskData.Description, taskData.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}

// DeleteTaskHandler удаляет задачу по ID
// DELETE /tasks/{id}
//
// Возвращает код 204 при успешном удалении
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage, id int) {
	err := storage.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
