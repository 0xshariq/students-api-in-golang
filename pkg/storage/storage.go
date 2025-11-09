package storage

import (
	"database/sql"

	"github.com/0xshariq/students-api-in-golang/pkg/types"
)

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	DeleteStudent(id int64) (sql.Result, error)
	UpdateStudent(id int64, name string, email string, age int) (sql.Result, error)
}