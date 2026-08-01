package user

import (
	"context"
	"expense-tracker-api/internal/httpserver"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo     *UserRepository
	Validate *validator.Validate
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		Repo:     repo,
		Validate: validator.New(),
	}
}

func (u *UserService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	if err := u.Validate.Struct(req); err != nil {
		return nil, httpserver.NewClientError(http.StatusBadRequest, formatValidationError(err))
	}

	existing, err := u.Repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, httpserver.NewClientError(http.StatusConflict, "email is already registered")
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return u.Repo.CreateUser(ctx, req.Name, req.Email, passwordHash)
}

func (u *UserService) Login(ctx context.Context, req LoginRequest) (*User, error) {
	invalidCredentials := httpserver.NewClientError(http.StatusUnauthorized, "invalid email or password")

	if err := u.Validate.Struct(req); err != nil {
		return nil, invalidCredentials
	}

	existing, err := u.Repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, invalidCredentials
	}

	if !comparePassword(req.Password, existing.PasswordHash) {
		return nil, invalidCredentials
	}

	return existing, nil
}

func formatValidationError(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok || len(validationErrors) == 0 {
		return "invalid request"
	}

	field := validationErrors[0]
	switch field.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field.Field(), field.Param())
	default:
		return fmt.Sprintf("%s is invalid", field.Field())
	}
}

func hashPassword(password string) (string, error) {
	passByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passByte), nil
}

func comparePassword(password, passwordHash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return false
	}
	return true
}
