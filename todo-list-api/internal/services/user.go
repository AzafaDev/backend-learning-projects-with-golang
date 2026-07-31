package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"todo-list-api/internal/config"
	"todo-list-api/internal/database"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	validator *validator.Validate
	cfg       *config.Config
}

func NewUserService(repo *repository.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		repo:      repo,
		validator: validator.New(),
		cfg:       cfg,
	}
}

func (s *UserService) Register(ctx context.Context, req models.RegisterRequest) (string, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Error().Err(err).Msg("register: validation failed")
		return "", fmt.Errorf("name, email and password are required")
	}
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return "", fmt.Errorf("register: %w", err)
	}
	if existingUser != nil {
		return "", fmt.Errorf("email is already exists")
	}
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Msg("register: failed to hash password")
		return "", fmt.Errorf("something went wrong")
	}
	createdUser, err := s.repo.CreateUser(ctx, req.Name, req.Email, passwordHash)
	if err != nil {
		return "", fmt.Errorf("register: %w", err)
	}
	signedToken, err := generateToken(createdUser.ID, createdUser.Name, createdUser.Email, s.cfg.JwtSecretKey)
	if err != nil {
		log.Error().Err(err).Msg("register: failed to generate token")
		return "", fmt.Errorf("something went wrong")
	}
	return signedToken, nil
}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	if err := s.validator.Struct(req); err != nil {
		log.Error().Err(err).Msg("login: validation failed")
		return "", fmt.Errorf("invalid email or password")
	}
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, database.ErrNotFound) {
			log.Error().Err(err).Msg("login: failed to fetch user")
		}
		return "", fmt.Errorf("invalid email or password")
	}
	if !comparePassword(existingUser.PasswordHash, req.Password) {
		return "", fmt.Errorf("invalid email or password")
	}
	signedToken, err := generateToken(existingUser.ID, existingUser.Name, existingUser.Email, s.cfg.JwtSecretKey)
	if err != nil {
		log.Error().Err(err).Msg("login: failed to generate token")
		return "", fmt.Errorf("something went wrong")
	}
	return signedToken, nil
}

func hashPassword(password string) (string, error) {
	passByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passByte), nil
}

func comparePassword(hash, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false
	}
	return true
}

func generateToken(id uuid.UUID, name, email, jwtSecretKey string) (string, error) {
	claims := models.UserClaims{
		ID:       id,
		Username: name,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString, secretKey string) (*models.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.UserClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algorithm method does not match")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token is invalid")
}
