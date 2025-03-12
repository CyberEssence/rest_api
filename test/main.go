package main

import (
	"fmt"
	"net/http"
	"test/handlers"
	"test/storage"
)

func main() {
	// Инициализация хранилища и обработчиков
	taskStorage := storage.NewInMemoryStorage()
	mux := handlers.SetupHandlers(taskStorage)

	fmt.Println("Сервер запущен на порту 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %v\n", err)
		return
	}
}
