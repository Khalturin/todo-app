// backend/db/db.go
package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	// Убедись, что папка data существует
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err = os.Mkdir("./data", 0755)
		if err != nil {
			log.Fatal("Не удалось создать папку ./data:", err)
		}
	}

	var err error
	DB, err = sql.Open("sqlite3", "./data/todo.db")
	if err != nil {
		log.Fatal("Не удалось открыть БД:", err)
	}

	// Проверим подключение
	if err = DB.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	// Создадим таблицу задач (если ещё не создана)
	createTasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		parent_id INTEGER,
		category TEXT,
		notes TEXT,
		completed BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTasksTable)
	if err != nil {
		log.Fatal("Не удалось создать таблицу tasks:", err)
	}

	log.Println("База данных инициализирована: ./data/todo.db")
}