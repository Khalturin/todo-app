// backend/handlers/tasks.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"todo-app/backend/db"
	"todo-app/backend/models"
)

func GetTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, parent_id, category, notes, completed FROM tasks ORDER BY created_at")
	if err != nil {
		log.Printf("Ошибка при выборке задач: %v", err)
		http.Error(w, "Внутренняя ошибка", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		var parentID sql.NullInt64
		err := rows.Scan(&task.ID, &task.Title, &parentID, &task.Category, &task.Notes, &task.Completed)
		if err != nil {
			log.Printf("Ошибка при чтении строки: %v", err)
			http.Error(w, "Ошибка данных", http.StatusInternalServerError)
			return
		}
		if parentID.Valid {
			pid := int(parentID.Int64)
			task.ParentID = &pid
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	var parentID interface{}
	if task.ParentID != nil {
		parentID = *task.ParentID
	} else {
		parentID = nil
	}

	result, err := db.DB.Exec(
		"INSERT INTO tasks (title, parent_id, category, notes, completed) VALUES (?, ?, ?, ?, ?)",
		task.Title,
		parentID,
		task.Category,
		task.Notes,
		task.Completed,
	)
	if err != nil {
		log.Printf("Ошибка при вставке задачи: %v", err)
		http.Error(w, "Не удалось создать задачу", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	task.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
