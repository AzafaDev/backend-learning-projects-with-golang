package handlers

import (
	"encoding/json"
	"net/http"
	"todo-list-api/internal/models"
	"todo-list-api/internal/services"

	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	srv *services.UserService
}

func NewUserHandler(srv *services.UserService) *UserHandler {
	return &UserHandler{
		srv: srv,
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, message, status)
}

func writeJSON(w http.ResponseWriter, body any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	tokenString, err := h.srv.Login(r.Context(), req)
	if err != nil {
		writeJSONError(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	writeJSON(w, map[string]string{"token": tokenString})
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	tokenString, err := h.srv.Register(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("register failed")
		writeJSONError(w, "unable to register user", http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]string{"token": tokenString})
}
