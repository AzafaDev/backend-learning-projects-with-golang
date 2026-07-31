package services

import (
	"context"
	"fmt"
	"time"
	"todo-list-api/internal/config"
	"todo-list-api/internal/models"
	"todo-list-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
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

func (u *UserService) Register(ctx context.Context, req models.RegisterRequest) (string, error) {
	if err := u.validator.Struct(req); err != nil {
		log.Error().Err(err)
		return "", fmt.Errorf("name, email and password are required")
	}
	existingUser, err := u.repo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return "", fmt.Errorf("email is already exists")
	} else if err != nil {
		return "", err
	}
	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		log.Error().Err(err)
		return "", fmt.Errorf("something went wrong")
	}
	createdUser, err := u.repo.CreateUser(ctx, req.Name, req.Email, passwordHash)
	if err != nil {
		return "", err
	}
	signedToken, err := generateToken(createdUser.Name, createdUser.Email, u.cfg.JwtSecretKey)
	if err != nil {
		log.Error().Err(err)
		return "", err
	}
	return signedToken, nil
}

func (u *UserService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	if err := u.validator.Struct(req); err != nil {
		log.Error().Err(err)
		return "", err
	}
	existingUser, err := u.repo.GetUserByEmail(ctx, req.Email)
	if existingUser == nil || err != nil {
		log.Error().Err(err)
		return "", fmt.Errorf("invalid email or password")
	}
	if matchingPassword := comparePassword(existingUser.PasswordHash, req.Password); !matchingPassword {
		log.Error().Err(err)
		return "", fmt.Errorf("invalid email or password")
	}
	signedToken, err := generateToken(existingUser.Name, existingUser.Email, u.cfg.JwtSecretKey)
	if err != nil {
		log.Error().Err(err)
		return "", err
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

func generateToken(name, email, jwtSecretKey string) (string, error) {
	claims := models.UserClaims{
		Username: name,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func validateToken(tokenString, secretKey string) (*models.UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, models.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algorithm method does not match")
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token is invalid")
}
