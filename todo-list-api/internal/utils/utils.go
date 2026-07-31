package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, message, status)
}

func WriteJSON(w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body == nil || status == http.StatusNoContent {
		return
	}
	json.NewEncoder(w).Encode(body)
}
