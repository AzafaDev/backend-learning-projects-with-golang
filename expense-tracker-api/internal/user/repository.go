package user

import (
	"context"
	"errors"
	"fmt"

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

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `
	SELECT id, name, email, password_hash, created_at, updated_at
	FROM users
	WHERE email=$1
	`
	err := u.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, email, passwordHash string) (*User, error) {
	var user User
	query := `
	INSER INTO users
	(email, password_hash)
	VALUES ($1, $2)
	RETURNING id, name, email, password_hash, created_at, updated_at
	`
	err := u.DB.QueryRow(ctx, query, email, passwordHash).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
