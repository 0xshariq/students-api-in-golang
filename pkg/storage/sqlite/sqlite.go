package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/0xshariq/students-api-in-golang/pkg/config"
	"github.com/0xshariq/students-api-in-golang/pkg/types"
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

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	statement, error := s.DB.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if error != nil {
		return 0, error
	}
	defer statement.Close()

	result, err := statement.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Sqlite) GetStudentByID(id int64) (types.Student, error) {

	statement, err := s.DB.Prepare("SELECT id,name,email,age FROM students WHERE id = ? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}
	defer statement.Close()

	var student types.Student
	err = statement.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil

}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {

	statement, err := s.DB.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// func (s *Sqlite) DeleteStudent(id int64) (string, error) {

// 	statement, err := s.DB.Prepare("DELETE FROM students WHERE id = ?")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer statement.Close()

// 	_, err = statement.Exec(id)
// 	if err != nil {
// 		return "", err
// 	}

// 	return "Student deleted successfully", nil
// }

// func (s *Sqlite) UpdateStudent(id int64, name string, email string, age int) (string, error) {

// 	statement, err := s.DB.Prepare("UPDATE students SET name = ?, email = ?, age = ? WHERE id = ?")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer statement.Close()

// 	_, err = statement.Exec(name, email, age, id)
// 	if err != nil {
// 		return "", err
// 	}

// 	return "Student updated successfully", nil
// }
