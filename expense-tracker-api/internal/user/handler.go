package user

import (
	"encoding/json"
	"expense-tracker-api/internal/config"
	"expense-tracker-api/internal/httpserver"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	srv *UserService
	cfg *config.Config
}

func NewUserHandler(srv *UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		srv: srv,
		cfg: cfg,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("register: user handler")
		httpserver.WriteJSONError(w, "name, email and password are required", http.StatusBadRequest)
		return
	}
	createdUser, err := u.srv.Register(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("register: user handler")
		httpserver.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := GenerateJwtToken(createdUser.ID, u.cfg.JWTSecretKey)
	if err != nil {
		log.Error().Err(err).Msg("register: user handler")
		httpserver.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpserver.WriteJSON(w, token, http.StatusCreated)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("login: user handler")
		httpserver.WriteJSONError(w, "email and password are required", http.StatusBadRequest)
		return
	}
	existingUser, err := u.srv.Login(r.Context(), req)
	if err != nil {
		log.Error().Err(err).Msg("login: user handler")
		httpserver.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := GenerateJwtToken(existingUser.ID, u.cfg.JWTSecretKey)
	if err != nil {
		log.Error().Err(err).Msg("login: user handler")
		httpserver.WriteJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	httpserver.WriteJSON(w, token, http.StatusCreated)
}

func GenerateJwtToken(id uuid.UUID, secret string) (string, error) {
	claims := UserClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString, secretKey string) (*UserClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is not valid")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("token is invalid")
}
