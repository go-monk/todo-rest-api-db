package todo

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteTaskStore struct {
	db *sql.DB
}

func NewSqliteTaskStore(dbPath string) (*SqliteTaskStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT NOT NULL
    )`)
	if err != nil {
		return nil, err
	}
	return &SqliteTaskStore{db: db}, nil
}

// Change the interface to return (int, error) for better error handling
func (s *SqliteTaskStore) CreateTask(text string) (int, error) {
	res, err := s.db.Exec("INSERT INTO tasks (text) VALUES (?)", text)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *SqliteTaskStore) GetTasks() []Task {
	rows, err := s.db.Query("SELECT id, text FROM tasks")
	if err != nil {
		return nil
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.Id, &t.Text); err == nil {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (s *SqliteTaskStore) GetTask(id int) (Task, error) {
	var t Task
	err := s.db.QueryRow("SELECT id, text FROM tasks WHERE id = ?", id).Scan(&t.Id, &t.Text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, errors.New("task not found")
		}
		return Task{}, err
	}
	return t, nil
}

func (s *SqliteTaskStore) DeleteTask(id int) error {
	res, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("task not found")
	}
	return nil
}
