package repository

import (
	"context"
	"errors"
	"fmt"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, name, email, password_hash, created_at, updated_at
	FROM users
	WHERE email=$1
	`
	err := r.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, database.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, name, email, passwordHash string) (*models.User, error) {
	var user models.User
	query := `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, name, email, password_hash, created_at, updated_at
	`
	err := r.DB.QueryRow(ctx, query, name, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}
