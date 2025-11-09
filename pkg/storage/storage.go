package storage

import "github.com/0xshariq/students-api-in-golang/pkg/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentByID(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	// DeleteStudent(id string) (string, error)
	// UpdateStudent(id string, name string, email string, age int) (string, error)
}