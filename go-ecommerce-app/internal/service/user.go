package service

import (
	"context"
	"errors"
	"fmt"
	"go-ecommerce-app/internal/config"
	"go-ecommerce-app/internal/model"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/pkg/security"

	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolationCode = "23505"

type UserService struct {
	repo *repository.Queries
	cfg  *config.Config
}

func NewUserService(repo *repository.Queries, cfg *config.Config) *UserService {
	return &UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (u *UserService) Register(ctx context.Context, req model.RegisterUserRequest) (*model.UserResponse, error) {
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := u.repo.CreateUser(ctx, repository.CreateUserParams{
		FullName:     req.FullName,
		Email:        req.Email,
		PasswordHash: passwordHash,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("error in creating user: %w", err)
	}

	return &model.UserResponse{
		ID:        createdUser.ID.Bytes,
		FullName:  createdUser.FullName,
		Email:     createdUser.Email,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt.Time,
		UpdatedAt: createdUser.UpdatedAt.Time,
	}, nil
}

func (u *UserService) Login(ctx context.Context, req model.LoginUserRequest) (*model.UserResponse, string, error) {
	existingUser, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("error in getting user by email: %w", err)
	}

	if err := security.ComparePassword(existingUser.PasswordHash, req.Password); err != nil {
		return nil, "", fmt.Errorf("error in comparing password: %w", err)
	}

	signedToken, err := security.GenerateToken(u.cfg.JWTSecret, u.cfg.JWTExpiry, model.UserResponse{
		ID:        existingUser.ID.Bytes,
		FullName:  existingUser.FullName,
		Email:     existingUser.Email,
		Role:      existingUser.Role,
		CreatedAt: existingUser.CreatedAt.Time,
		UpdatedAt: existingUser.UpdatedAt.Time,
	})

	if err != nil {
		return nil, "", fmt.Errorf("error in generating token JWT: %w", err)
	}

	return &model.UserResponse{
		ID:        existingUser.ID.Bytes,
		FullName:  existingUser.FullName,
		Email:     existingUser.Email,
		Role:      existingUser.Role,
		CreatedAt: existingUser.CreatedAt.Time,
		UpdatedAt: existingUser.UpdatedAt.Time,
	}, signedToken, nil
}
