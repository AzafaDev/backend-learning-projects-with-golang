package services

import (
	"context"
	"errors"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TodoService struct {
	Repo      *repository.TodoRepository
	validator *validator.Validate
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{
		Repo:      repo,
		validator: validator.New(),
	}
}

func (t *TodoService) GetTodos(ctx context.Context, userID uuid.UUID) ([]models.Todo, error) {
	return t.Repo.GetTodosByUserID(ctx, userID)
}

func (t *TodoService) CreateTodo(ctx context.Context, req models.CreateTodoRequest, userID uuid.UUID) (*models.Todo, error) {
	if err := t.validator.Struct(req); err != nil {
		return nil, errors.New("title and description are required")
	}
	return t.Repo.CreateTodo(ctx, req, userID)
}

func (t *TodoService) UpdateTodo(ctx context.Context, req models.UpdateTodoRequest, id, userID uuid.UUID) (*models.Todo, error) {
	if err := t.validator.Struct(req); err != nil {
		return nil, errors.New("title and description are required")
	}
	return t.Repo.UpdateTodo(ctx, req, id, userID)
}

func (t *TodoService) DeleteTodo(ctx context.Context, id, userID uuid.UUID) (*models.Todo, error) {
	return t.Repo.DeleteTodo(ctx, id, userID)
}
