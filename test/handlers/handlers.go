// handlers/handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"test/storage"
)

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

func CreateTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage) {
	var taskData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if taskData.Title == "" || taskData.Description == "" {
		http.Error(w, "Title и Description обязательны", http.StatusBadRequest)
		return
	}

	task, err := storage.CreateTask(taskData.Title, taskData.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetAllTasksHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage) {
	tasks, err := storage.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage, id int) {
	task, err := storage.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage, id int) {
	var taskData struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	err := json.NewDecoder(r.Body).Decode(&taskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := storage.UpdateTask(id, taskData.Title, taskData.Description, taskData.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request, storage *storage.InMemoryStorage, id int) {
	err := storage.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
