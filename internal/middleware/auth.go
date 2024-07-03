package middleware

import (
	"budgeting-app/internal/auth"
	"context"
	"net/http"
	"strings"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// You can pass claims to context if needed
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		http.Error(w, claims.Username, 123)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
