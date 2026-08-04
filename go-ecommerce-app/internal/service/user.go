package service

import (
	"context"
	"errors"
	"fmt"
	"go-ecommerce-app/internal/config"
	"go-ecommerce-app/internal/model"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/pkg/security"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

var expiryVerificationEmail = time.Now().Add(24 * time.Hour)

const pgUniqueViolationCode = "23505"

type UserService struct {
	repo  *repository.Queries
	cfg   *config.Config
	email EmailSender
}

func NewUserService(repo *repository.Queries, cfg *config.Config, email EmailSender) *UserService {
	return &UserService{
		repo:  repo,
		cfg:   cfg,
		email: email,
	}
}

func (u *UserService) Register(ctx context.Context, req model.RegisterUserRequest) (*model.UserResponse, error) {
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	trimEmail := strings.ToLower(strings.TrimSpace(req.Email))

	createdUser, err := u.repo.CreateUser(ctx, repository.CreateUserParams{
		FullName:     req.FullName,
		Email:        trimEmail,
		PasswordHash: passwordHash,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolationCode {
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("error in creating user: %w", err)
	}

	randomString, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("error in generatin random string token: %w", err)
	}
	hashedToken := security.HashRefreshToken(randomString)

	_, err = u.repo.CreateVericationEmail(ctx, repository.CreateVericationEmailParams{
		UserID:    createdUser.ID,
		TokenHash: hashedToken,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiryVerificationEmail,
			Valid: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error in creating verification email token: %w", err)
	}

	if err := u.email.SendVerificationEmail(ctx, createdUser.Email, randomString); err != nil {
		return nil, fmt.Errorf("error in sending verification email: %w", err)
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

func (u *UserService) Login(ctx context.Context, req model.LoginUserRequest) (*model.UserResponse, string, string, error) {
	existingUser, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", "", fmt.Errorf("error in getting user by email: %w", err)
	}

	if err := security.ComparePassword(existingUser.PasswordHash, req.Password); err != nil {
		return nil, "", "", fmt.Errorf("error in comparing password: %w", err)
	}

	signedToken, err := security.GenerateToken(u.cfg.JWTSecret, u.cfg.JWTExpiry, existingUser.ID.Bytes, existingUser.Role)

	if err != nil {
		return nil, "", "", fmt.Errorf("error in generating token JWT: %w", err)
	}

	rawToken, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, "", "", fmt.Errorf("error in generating refresh token: %w", err)
	}
	hashedToken := security.HashRefreshToken(rawToken)

	if err := u.repo.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:    existingUser.ID,
		TokenHash: hashedToken,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(u.cfg.RefreshTokenExpiry),
			Valid: true,
		},
	}); err != nil {
		return nil, "", "", fmt.Errorf("error in creating refresh token: %w", err)
	}

	return &model.UserResponse{
		ID:        existingUser.ID.Bytes,
		FullName:  existingUser.FullName,
		Email:     existingUser.Email,
		Role:      existingUser.Role,
		CreatedAt: existingUser.CreatedAt.Time,
		UpdatedAt: existingUser.UpdatedAt.Time,
	}, signedToken, rawToken, nil
}

func (u *UserService) Refresh(ctx context.Context, rawToken string) (string, string, error) {
	tokenHash := security.HashRefreshToken(rawToken)

	existingRefreshToken, err := u.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return "", "", fmt.Errorf("error in get refresh token: %w", err)
	}

	if err := u.repo.RevokeRefreshToken(ctx, existingRefreshToken.TokenHash); err != nil {
		return "", "", fmt.Errorf("error in get refresh token: %w", err)
	}

	newRawToken, err := security.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("error in generating refresh token: %w", err)
	}
	hashedToken := security.HashRefreshToken(newRawToken)

	if err := u.repo.CreateRefreshToken(ctx, repository.CreateRefreshTokenParams{
		UserID:    existingRefreshToken.UserID,
		TokenHash: hashedToken,
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(u.cfg.RefreshTokenExpiry),
			Valid: true,
		},
	}); err != nil {
		return "", "", fmt.Errorf("error in creating refresh token: %w", err)
	}

	existingUser, err := u.repo.GetUserByID(ctx, existingRefreshToken.UserID)
	if err != nil {
		return "", "", fmt.Errorf("error in getting user by id: %w", err)
	}

	accessToken, err := security.GenerateToken(u.cfg.JWTSecret, u.cfg.JWTExpiry, existingRefreshToken.UserID.Bytes, existingUser.Role)
	if err != nil {
		return "", "", fmt.Errorf("error in generating token: %w", err)
	}

	return accessToken, newRawToken, nil
}

func (u *UserService) Logout(ctx context.Context, rawToken string) error {
	tokenHash := security.HashRefreshToken(rawToken)

	if err := u.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("error in revoking refresh token: %w", err)
	}

	return nil
}
