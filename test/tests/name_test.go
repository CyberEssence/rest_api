// Package tests содержит тесты для API обработчиков задач
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"test/handlers"
	"test/models"
	"test/storage"
	"testing"
)

// TestCreateTaskHandler проверяет создание новой задачи через POST /tasks
//
// Проверяет:
// - Корректность создания задачи
// - Правильность установки ID
// - Соответствие полей созданной задачи
func TestCreateTaskHandler(t *testing.T) {
	// Инициализация хранилища и обработчиков
	taskStorage := storage.NewInMemoryStorage()
	mux := handlers.SetupHandlers(taskStorage)

	// Эталонная задача для тестирования
	expectedTask := models.Task{
		Title:       "Купить продукты",
		Description: "Молоко, хлеб, овощи",
	}

	// Преобразование эталонной задачи в JSON
	jsonTask, err := json.Marshal(expectedTask)
	if err != nil {
		t.Fatal(err)
	}

	// Создание POST запроса
	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	if err != nil {
		t.Fatal(err)
	}

	// Запись ответа
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Проверка статуса ответа
	if w.Code != http.StatusCreated {
		t.Errorf("Ожидался код %d, получен %d", http.StatusCreated, w.Code)
	}

	// Получение тела ответа
	body := w.Body.Bytes()
	t.Logf("Ответ сервера: %s", body)

	// Десериализация полученной задачи
	var actualTask models.Task
	err = json.Unmarshal(body, &actualTask)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка соответствия всех полей
	if actualTask.Title != expectedTask.Title {
		t.Errorf("Несовпадение заголовка:\nОжидалось: %q\nПолучено: %q", expectedTask.Title, actualTask.Title)
	}
	if actualTask.Description != expectedTask.Description {
		t.Errorf("Несовпадение описания:\nОжидалось: %q\nПолучено: %q", expectedTask.Description, actualTask.Description)
	}
	if actualTask.ID == 0 {
		t.Errorf("ID задачи не был установлен")
	}
}

// TestGetAllTasksHandler проверяет получение списка всех задач через GET /tasks
//
// Проверяет:
// - Корректность получения списка задач
// - Соответствие количества задач
// - Правильность данных возвращаемых задач
func TestGetAllTasksHandler(t *testing.T) {
	// Инициализация хранилища и обработчиков
	taskStorage := storage.NewInMemoryStorage()
	mux := handlers.SetupHandlers(taskStorage)

	// Создание тестовой задачи
	testTask := models.Task{
		Title:       "Тестовая задача",
		Description: "Описание тестовой задачи",
	}
	taskStorage.CreateTask(testTask.Title, testTask.Description)

	// Создание GET запроса
	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Запись ответа
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Проверка статуса ответа
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался код %d, получен %d", http.StatusOK, w.Code)
	}

	// Получение тела ответа
	body := w.Body.Bytes()
	t.Logf("Ответ сервера: %s", body)

	// Десериализация полученных задач
	var tasks []*models.Task
	err = json.Unmarshal(body, &tasks)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка количества задач
	if len(tasks) != 1 {
		t.Errorf("Ожидалось 1 задача, получено %d", len(tasks))
	}

	// Проверка соответствия полей
	if tasks[0].Title != testTask.Title || tasks[0].Description != testTask.Description {
		t.Errorf("Несовпадение данных задачи")
	}
}

// TestGetTaskHandler проверяет получение задачи по ID через GET /tasks/{id}
//
// Проверяет:
// - Корректность получения задачи по ID
// - Соответствие всех полей задачи
func TestGetTaskHandler(t *testing.T) {
	// Инициализация хранилища и обработчиков
	taskStorage := storage.NewInMemoryStorage()
	mux := handlers.SetupHandlers(taskStorage)

	// Создание тестовой задачи
	testTask := models.Task{
		Title:       "Тестовая задача",
		Description: "Описание тестовой задачи",
	}
	taskStorage.CreateTask(testTask.Title, testTask.Description)

	// Создание GET запроса
	req, err := http.NewRequest("GET", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Запись ответа
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Проверка статуса ответа
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался код %d, получен %d", http.StatusOK, w.Code)
	}

	// Получение тела ответа
	body := w.Body.Bytes()
	t.Logf("Ответ сервера: %s", body)

	// Десериализация полученной задачи
	var task models.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка соответствия полей
	if task.Title != testTask.Title || task.Description != testTask.Description {
		t.Errorf("Несовпадение данных задачи")
	}
}

// TestUpdateTaskHandler проверяет обновление задачи через PUT /tasks/{id}
//
// Проверяет:
// - Корректность обновления всех полей задачи
// - Соответствие обновленных данных
func TestUpdateTaskHandler(t *testing.T) {
	// Инициализация хранилища и обработчиков
	taskStorage := storage.NewInMemoryStorage()
	mux := handlers.SetupHandlers(taskStorage)

	// Создание тестовой задачи
	testTask := models.Task{
		Title:       "Тестовая задача",
		Description: "Описание тестовой задачи",
	}
	taskStorage.CreateTask(testTask.Title, testTask.Description)

	// Обновление задачи
	updatedTask := models.Task{
		Title:       "Обновленная задача",
		Description: "Новое описание",
		Completed:   true,
	}

	// Преобразование обновленной задачи в JSON
	jsonTask, err := json.Marshal(updatedTask)
	if err != nil {
		t.Fatal(err)
	}

	// Создание PUT запроса
	req, err := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(jsonTask))
	if err != nil {
		t.Fatal(err)
	}

	// Запись ответа
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Проверка статуса ответа
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался код %d, получен %d", http.StatusOK, w.Code)
	}

	// Получение тела ответа
	body := w.Body.Bytes()
	t.Logf("Ответ сервера: %s", body)

	// Десериализация обновленной задачи
	var task models.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		t.Fatal(err)
	}

	// Проверка соответствия полей
	if task.Title != updatedTask.Title || task.Description != updatedTask.Description || task.Completed != updatedTask.Completed {
		t.Errorf("Несовпадение данных обновленной задачи")
	}
}

// TestDeleteTaskHandler проверяет удаление задачи через DELETE /tasks/{id}
//
// Проверяет:
// - Корректность удаления задачи
// - Отсутствие задачи после удаления
func TestDeleteTaskHandler(t *testing.T) {
	// Инициализация хранилища и обработчиков
	taskStorage := storage.NewInMemoryStorage()
	mux := handlers.SetupHandlers(taskStorage)

	// Создание тестовой задачи
	testTask := models.Task{
		Title:       "Тестовая задача",
		Description: "Описание тестовой задачи",
	}
	taskStorage.CreateTask(testTask.Title, testTask.Description)

	// Создание DELETE запроса
	req, err := http.NewRequest("DELETE", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Запись ответа
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	// Проверка статуса ответа
	if w.Code != http.StatusNoContent {
		t.Errorf("Ожидался код %d, получен %d", http.StatusNoContent, w.Code)
	}

	// Проверка, что задача удалена
	req, err = http.NewRequest("GET", "/tasks/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Ожидался код %d, получен %d", http.StatusNotFound, w.Code)
	}
}
