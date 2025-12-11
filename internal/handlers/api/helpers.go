package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RespondJSON отправляет JSON ответ
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// RespondError отправляет ошибку в формате JSON
func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, ErrorResponse{
		Code:    status,
		Message: message,
	})
}

// DecodeJSON парсит JSON из request body
func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// Error handlers
func RespondUserNotFound(w http.ResponseWriter) {
	RespondError(w, http.StatusNotFound, "User not found")
}

func RespondUserAlreadyExists(w http.ResponseWriter) {
	RespondError(w, http.StatusConflict, "User with this Telegram ID already exists")
}

func RespondCarNotFound(w http.ResponseWriter) {
	RespondError(w, http.StatusNotFound, "Car not found")
}

func RespondCarAccessDenied(w http.ResponseWriter) {
	RespondError(w, http.StatusForbidden, "Access denied to this car")
}

func RespondBadRequest(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusBadRequest, message)
}

func RespondUnauthorized(w http.ResponseWriter, message string) {
	RespondError(w, http.StatusUnauthorized, message)
}

func RespondInternalError(w http.ResponseWriter) {
	RespondError(w, http.StatusInternalServerError, "Internal server error")
}
