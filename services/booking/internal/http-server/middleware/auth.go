package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	appjwt "github.com/antongolenev23/voltake-services/pkg/jwt"
)

var ErrClaimsNotFound = errors.New("claims not found")

type contextKey string

const ClaimsKey contextKey = "claims"

func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			if tokenString == authHeader {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			claims := &appjwt.Claims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(token *jwt.Token) (any, error) {
					if token.Method != jwt.SigningMethodHS256 {
						return nil, jwt.ErrSignatureInvalid
					}

					return []byte(secret), nil
				},
			)

			if err != nil || !token.Valid {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				ClaimsKey,
				claims,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) (*appjwt.Claims, error) {
	claims, ok := ctx.Value(ClaimsKey).(*appjwt.Claims)

	if !ok {
		return nil, ErrClaimsNotFound
	}

	return claims, nil
}

func IsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, err := GetClaims(r.Context())

		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !claims.IsAdmin {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
