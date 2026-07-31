package handlers

import (
	"encoding/json"
	"net/http"
	"todo-list-api/internal/models"
	"todo-list-api/internal/services"
	"todo-list-api/internal/utils"

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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	tokenString, err := h.srv.Login(r.Context(), req)
	if err != nil {
		utils.WriteJSONError(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	utils.WriteJSON(w, map[string]string{"token": tokenString}, http.StatusOK)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	tokenString, err := h.srv.Register(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("register failed")
		utils.WriteJSONError(w, "unable to register user", http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, map[string]string{"token": tokenString}, http.StatusCreated)
}
