package repository

import (
	"context"
	"errors"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	Db *pgxpool.Pool
}

func NewTodoRepository(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		Db: db,
	}
}

func (t *TodoRepository) GetTodoByID(ctx context.Context, id uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	query := `
	SELECT id, user_id, title, description, created_at, updated_at
	FROM todos
	WHERE id=$1
	`
	if err := t.Db.QueryRow(ctx, query, id).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &todo, nil
}

func (t *TodoRepository) CreateTodo(ctx context.Context, req models.CreateTodoRequest, userID uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	query := `
	INSERT INTO users (title, description, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, user_id, title, description, created_at, updated_at
	`
	if err := t.Db.QueryRow(ctx, query, req.Title, req.Description, userID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (t *TodoRepository) UpdateTodo(ctx context.Context, req models.UpdateTodoRequest, id, userID uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	query := `
	UPDATE todos
	SET title=$1
		description=$2
		updated_at=now()
	WHERE id=$3
	AND user_id=$4
	RETURNING id, user_id, title, description, created_at, updated_at
	`
	if err := t.Db.QueryRow(ctx, query, req.Title, req.Description, id, userID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &todo, nil
}

func (t *TodoRepository) DeleteTodo(ctx context.Context, id, userID uuid.UUID) (*models.Todo, error) {
	var todo models.Todo
	query := `
	DELETE todos
	WHERE id=$1
	AND user_id=$2
	RETURNING id, user_id, title, description, created_at, updated_at
	`
	if err := t.Db.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, err
	}

	return &todo, nil
}

func (t *TodoRepository) GetTodosByUserID(ctx context.Context, userID uuid.UUID) ([]models.Todo, error) {
	var todos []models.Todo
	query := `
	SELECT id, user_id, title, description, created_at, updated_at
	FROM todos
	WHERE user_id=$1
	`

	rows, err := t.Db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Description,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if todos == nil {
		return []models.Todo{}, nil
	}
	return todos, nil
}
