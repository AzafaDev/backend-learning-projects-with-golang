package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIErrorRespose struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func WriteError(w http.ResponseWriter, status int, message string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	resp := APIErrorRespose{
		Success: false,
		Message: message,
		Error:   errMsg,
	}
	WriteJSON(w, status, resp)
}
