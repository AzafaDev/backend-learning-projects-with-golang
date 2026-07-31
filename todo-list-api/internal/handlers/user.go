package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"todo-list-api/internal/models"
	"todo-list-api/internal/services"
)

type UserHandler struct {
	srv *services.UserService
}

func NewUserHandler(srv *services.UserService) *UserHandler {
	return &UserHandler{
		srv: srv,
	}
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "invalid or email password", http.StatusBadRequest)
		return
	}
	tokenString, err := u.srv.Login(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "invalid or email password", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokenString); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "invalid or email password", http.StatusBadRequest)
		return
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "name, email and password are required", http.StatusBadRequest)
		return
	}
	tokenString, err := u.srv.Register(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokenString); err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
