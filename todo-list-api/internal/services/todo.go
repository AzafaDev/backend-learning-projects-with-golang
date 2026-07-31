package services

import (
	"context"
	"errors"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ErrForbidden = errors.New("forbidden")

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

func (t *TodoService) GetTodos(ctx context.Context, userID uuid.UUID, page, limit int) (*models.PaginatedTodos, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	todos, err := t.Repo.GetTodosByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	total, err := t.Repo.CountTodosByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.PaginatedTodos{
		Data:  todos,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
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
	existing, err := t.Repo.GetTodoByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
		return nil, ErrForbidden
	}
	return t.Repo.UpdateTodo(ctx, req, id, userID)
}

func (t *TodoService) DeleteTodo(ctx context.Context, id, userID uuid.UUID) (*models.Todo, error) {
	existing, err := t.Repo.GetTodoByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
		return nil, ErrForbidden
	}
	return t.Repo.DeleteTodo(ctx, id, userID)
}

func (t *TodoService) GetTodoByID(ctx context.Context, id, userID uuid.UUID) (*models.Todo, error) {
	return t.Repo.GetTodoByID(ctx, id)
}
