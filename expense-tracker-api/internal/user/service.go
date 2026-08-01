package user

import (
	"context"
	"fmt"

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
		return nil, err
	}

	existing, err := u.Repo.GetUserByEmail(ctx, req.Email)
	if existing != nil {
		return nil, fmt.Errorf("email is already exists")
	} else if err != nil {
		return nil, err
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	return u.Repo.CreateUser(ctx, req.Email, passwordHash)
}

func (u *UserService) Login(ctx context.Context, req LoginRequest) (*User, error) {
	if err := u.Validate.Struct(req); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	existing, err := u.Repo.GetUserByEmail(ctx, req.Email)
	if existing == nil {
		return nil, fmt.Errorf("invalid email or password")
	} else if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	matchingPassword := comparePassword(req.Password, existing.PasswordHash)
	if !matchingPassword {
		return nil, fmt.Errorf("invalid email or password")
	}

	return existing, nil
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
