package middleware

import (
	"context"
	"expense-tracker-api/internal/config"
	"expense-tracker-api/internal/httpserver"
	"expense-tracker-api/internal/user"
	"net/http"
	"strings"
)

type contextKey string

const userIdKey contextKey = "user_id"
const bearerPrefix = "Bearer "

func AuthMiddleware(next http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpserver.RespondError(w, "auth middleware", httpserver.NewClientError(http.StatusUnauthorized, "unauthorized"))
			return
		}

		tokenString := authHeader
		if ok := strings.HasPrefix(tokenString, bearerPrefix); ok {
			tokenString = strings.TrimPrefix(tokenString, bearerPrefix)
		}

		userClaims, err := user.ValidateToken(tokenString, cfg.JWTSecretKey)
		if err != nil {
			httpserver.RespondError(w, "auth middleware", httpserver.NewClientError(http.StatusUnauthorized, "token is invalid"))
			return
		}

		ctx := context.WithValue(r.Context(), userIdKey, userClaims.ID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
