package response

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]any

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, code, message string, details []FieldError) {
	JSON(w, status, Envelope{
		"error": ErrorResponse{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}
