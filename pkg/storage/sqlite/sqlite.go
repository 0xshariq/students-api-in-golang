package sqlite

import (
	"database/sql"

	"github.com/0xshariq/students-api-in-golang/pkg/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	DB *sql.DB
}

// New creates and returns a new Sqlite instance connected to the file
func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	// create table if not exists
	_, execErr := db.Exec(`CREATE TABLE IF NOT EXISTS students (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  email TEXT,
  age INTEGER
)`)
	if execErr != nil {
		return nil, execErr
	}

	return &Sqlite{DB: db}, nil
}
