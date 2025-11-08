package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/0xshariq/students-api-in-golang/pkg/types"
	"github.com/0xshariq/students-api-in-golang/pkg/utils/response"
	"github.com/go-playground/validator/v10"
)

func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to students api in golang"))
	}
}

func NewStudent() http.HandlerFunc {
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

		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "ok"})
	}
}
func CreateStudent(name string, email string, age int) (int64, error) {
	return 0, nil
}
func DeleteStudent(id int64) (string, error) {
	return "", nil
}
func UpdateStudent(id int64, data any) (string, error) {
	return "", nil
}