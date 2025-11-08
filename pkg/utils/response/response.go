package response

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func GeneralError(err error) ErrorResponse {
	return ErrorResponse{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationError(error validator.ValidationErrors) ErrorResponse {
	var errKeys []string
	for _, err := range error {
		switch err.Tag() {
		case "required":
			errKeys = append(errKeys, err.Field()+" is required")
		default:
			errKeys = append(errKeys, err.Field()+" is invalid")
		}
	}
	return ErrorResponse{
		Status: StatusError,
		Error: strings.Join(errKeys, ","),
	}
}
