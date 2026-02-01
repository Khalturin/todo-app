// backend/models/task.go
package models

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	ParentID  *int   `json:"parent_id,omitempty"` // nil = корневая задача
	Category  string `json:"category,omitempty"`
	Notes     string `json:"notes,omitempty"`
	Completed bool   `json:"completed"`
}
