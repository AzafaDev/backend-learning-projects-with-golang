package middleware

import (
	"context"
	"net/http"
	"strings"
	"todo-list-api/internal/config"
	"todo-list-api/internal/services"
	"todo-list-api/internal/utils"

	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "userID"
const bearerPrefix = "Bearer "

func AuthMiddleware(next http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.WriteJSONError(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		tokenString := authHeader
		if ok := strings.HasPrefix(authHeader, "Bearer"); ok {
			tokenString = strings.TrimPrefix(tokenString, bearerPrefix)
		}

		userClaims, err := services.ValidateToken(tokenString, cfg.JwtSecretKey)
		if err != nil {
			utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userClaims.ID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, false
	}
	return userID, true
}
