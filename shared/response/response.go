package response

import (
	"encoding/json"
	"net/http"
)

// ErrorBody is the standard API error envelope.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail holds a machine-readable code and human message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// DataBody wraps a successful payload.
type DataBody struct {
	Data any `json:"data"`
}

// JSON writes a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// Error writes a standard error envelope.
func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorBody{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}

// Data writes a standard success envelope.
func Data(w http.ResponseWriter, status int, payload any) {
	JSON(w, status, DataBody{Data: payload})
}
