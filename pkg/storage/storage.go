package storage

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	DeleteStudent(id int64) (string, error)
	UpdateStudent(id int64, data any) (string, error)
}