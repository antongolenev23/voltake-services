package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey  contextKey = "userID"
	IsAdminKey contextKey = "isAdmin"
)

func JWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(
					w,
					"missing authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			if tokenString == authHeader {
				http.Error(
					w,
					"invalid authorization format",
					http.StatusUnauthorized,
				)
				return
			}

			token, err := jwt.Parse(
				tokenString,
				func(token *jwt.Token) (any, error) {

					if token.Method != jwt.SigningMethodHS256 {
						return nil, jwt.ErrSignatureInvalid
					}

					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(
					w,
					"invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(
					w,
					"invalid claims",
					http.StatusUnauthorized,
				)
				return
			}

			userID, ok := claims["user_id"].(string)
			if !ok || userID == "" {
				http.Error(
					w,
					"missing user_id",
					http.StatusUnauthorized,
				)
				return
			}

			isAdmin, ok := claims["is_admin"].(bool)
			if !ok {
				http.Error(
					w,
					"missing is_admin",
					http.StatusUnauthorized,
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				userID,
			)

			ctx = context.WithValue(
				ctx,
				IsAdminKey,
				isAdmin,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}
