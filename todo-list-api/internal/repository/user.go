package repository

import (
	"context"
	"errors"
	"log"
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

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, name, email, password_hash, created_at, updated_at
	FROM users
	WHERE email=$1
	`
	if err := u.DB.QueryRow(ctx, query, email).Scan(&user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println(err)
			return nil, database.ErrNotFound("user with email: " + email)
		}
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) CreateUser(ctx context.Context, name, email, passwordHash string) (*models.User, error) {
	var user models.User
	query := `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING *
	`
	if err := u.DB.QueryRow(ctx, query, email).Scan(&user); err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
