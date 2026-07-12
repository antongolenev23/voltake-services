package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/antongolenev23/voltake-services/pkg/types"
)

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			id := uuid.NewString()

			ctx := context.WithValue(
				r.Context(),
				types.RequestIDKey,
				id,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
