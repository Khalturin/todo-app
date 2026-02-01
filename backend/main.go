// backend/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"todo-app/backend/db"
	"todo-app/backend/handlers"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "TODO API is running!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Инициализация БД
	db.InitDB()

	// API маршруты
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTasks(w, r)
		case http.MethodPost:
			handlers.CreateTask(w, r)
		default:
			http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		}
	})

	// Отдача статики: все остальные запросы — index.html или файлы из frontend/
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", http.StripPrefix("/", fs))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
