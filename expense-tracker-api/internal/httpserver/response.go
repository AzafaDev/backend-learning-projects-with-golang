package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ClientError struct {
	Status  int
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}

func NewClientError(status int, message string) *ClientError {
	return &ClientError{Status: status, Message: message}
}

func WriteError(w http.ResponseWriter, err error) {
	var clientErr *ClientError
	if errors.As(err, &clientErr) {
		WriteJSONError(w, clientErr.Message, clientErr.Status)
		return
	}

	WriteJSONError(w, "something went wrong, please try again later", http.StatusInternalServerError)
}

func WriteJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
