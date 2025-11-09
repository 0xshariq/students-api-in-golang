package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/0xshariq/students-api-in-golang/pkg/storage"
	"github.com/0xshariq/students-api-in-golang/pkg/types"
	"github.com/0xshariq/students-api-in-golang/pkg/utils/response"
	"github.com/go-playground/validator/v10"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api in golang"))
	}
}

func NewStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")
		var student types.Student
		error := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(error, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("empty body")))
			return
		}

		if error != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(error))
			return
		}

		// request validator
		if err := validator.New().Struct(student); err != nil {
			validateError := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateError))
			return
		}

		Id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Student created successfully", slog.Int64("id", Id))

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": Id})
	}
}

func GetStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Getting a student...", slog.String("id", id))

		// converting string id to int64
		intId, e := strconv.ParseInt(id, 10, 64)
		if e != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("invalid id format")))
			return
		}

		// Call the storage layer to get the student
		student, err := storage.GetStudentByID(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)

	}
}

func GetStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Getting all students...")

		// Call the storage layer to get all students
		students, err := storage.GetAllStudents()
		if err != nil {
			slog.Error("error getting students")
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, students)

	}
}

// func DeleteStudent(storage storage.Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := r.PathValue("id")
// 		slog.Info("Deleting a student...", slog.String("id", id))

// 		// Call the storage layer to delete the student
// 		message, err := storage.DeleteStudent(id)
// 		if err != nil {
// 			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}

// 		response.WriteJSON(w, http.StatusOK, map[string]string{"message": message})
// 	}
// }
// func UpdateStudent(storage storage.Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := r.PathValue("id")
// 		slog.Info("Updating a student...", slog.String("id", id))

// 		var student types.Student
// 		error := json.NewDecoder(r.Body).Decode(&student)
// 		if errors.Is(error, io.EOF) {
// 			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("empty body")))
// 			return
// 		}

// 		if error != nil {
// 			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(error))
// 			return
// 		}

// 		// request validator
// 		if err := validator.New().Struct(student); err != nil {
// 			validateError := err.(validator.ValidationErrors)
// 			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateError))
// 			return
// 		}

// 		// Call the storage layer to update the student
// 		message, err := storage.UpdateStudent(id, student.Name, student.Email, student.Age)
// 		if err != nil {
// 			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}

// 		response.WriteJSON(w, http.StatusOK, map[string]string{"message": message})
// 	}
// }
