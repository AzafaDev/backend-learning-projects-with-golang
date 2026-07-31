package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}
