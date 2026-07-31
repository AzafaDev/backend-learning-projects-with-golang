package handlers

import (
	"encoding/json"
	"net/http"
	"todo-list-api/internal/models"
	"todo-list-api/internal/services"
	"todo-list-api/internal/utils"

	"github.com/google/uuid"
)

type TodoHandler struct {
	srv *services.TodoService
}

func NewTodoHandler(srv *services.TodoService) *TodoHandler {
	return &TodoHandler{
		srv: srv,
	}
}

func (t *TodoHandler) CreateTodos(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		utils.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, "title and description are required", http.StatusBadRequest)
		return
	}

	createdTodo, err := t.srv.CreateTodo(r.Context(), req, userID)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, createdTodo, http.StatusCreated)
}

func (t *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateTodoRequest
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "id is invalid", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		utils.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, "title and description are required", http.StatusBadRequest)
		return
	}

	updatedTodo, err := t.srv.UpdateTodo(r.Context(), req, id, userID)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, updatedTodo, http.StatusCreated)
}

func (t *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		utils.WriteJSONError(w, "id is invalid", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		utils.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	_, err = t.srv.DeleteTodo(r.Context(), id, userID)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, nil, http.StatusNoContent)
}

func (t *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		utils.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	todos, err := t.srv.GetTodos(r.Context(), userID)
	if err != nil {
		utils.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WriteJSON(w, todos, http.StatusOK)
}
