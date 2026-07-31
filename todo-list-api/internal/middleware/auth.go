package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"todo-list-api/internal/config"
	"todo-list-api/internal/services"
	"todo-list-api/internal/utils"
)

func AuthMiddleware(next http.Handler, cfg config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"message": "Unauthorized",
			}); err != nil {
				http.Error(w, "something went wrong", http.StatusInternalServerError)
				return
			}
			return
		}
		tokenString := authHeader
		if ok := strings.HasPrefix(authHeader, "Bearer"); ok {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		userClaims, err := services.ValidateToken(tokenString, cfg.JwtSecretKey)
		if err != nil {
			utils.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userClaims.ID)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
