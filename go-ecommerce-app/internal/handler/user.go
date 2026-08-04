package handler

import (
	"encoding/json"
	"errors"
	"go-ecommerce-app/internal/model"
	"go-ecommerce-app/internal/service"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

type UserHandler struct {
	srv      *service.UserService
	validate *validator.Validate
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{
		srv:      srv,
		validate: validator.New(),
	}
}

func (u *UserHandler) UserRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Post("/login", u.Login)
		r.Post("/register", u.Register)
	})
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("register", "error", err)
		writeErrorJSON("invalid request payload", http.StatusBadRequest, w)
		return
	}

	if err := u.validate.Struct(req); err != nil {
		slog.Error("register", "error", err)
		writeErrorJSON("invalid request payload", http.StatusBadRequest, w)
		return
	}

	createdUser, err := u.srv.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailTaken) {
			slog.Error("register", "error", err)
			writeErrorJSON(err.Error(), http.StatusConflict, w)
			return
		}
		slog.Error("register", "error", err)
		writeErrorJSON("something went wrong", http.StatusInternalServerError, w)
		return
	}

	WriteJSON(JSONResponse{
		Success: true,
		Data: map[string]any{
			"user": createdUser,
		},
	}, http.StatusCreated, w)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("login", "error", err)
		writeErrorJSON("invalid request payload", http.StatusBadRequest, w)
		return
	}

	if err := u.validate.Struct(req); err != nil {
		slog.Error("login", "error", err)
		writeErrorJSON("invalid request payload", http.StatusBadRequest, w)
		return
	}

	existingUser, accessToken, err := u.srv.Login(r.Context(), req)
	if err != nil {
		slog.Error("login", "error", err)
		writeErrorJSON("invalid email or password", http.StatusBadRequest, w)
		return
	}

	WriteJSON(JSONResponse{
		Success: true,
		Data: map[string]any{
			"accessToken": accessToken,
			"user":        existingUser,
		},
	}, http.StatusOK, w)
}
