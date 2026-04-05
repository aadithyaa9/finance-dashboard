package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/aadithyaa9/finance-dashboard/internal/auth"
	"github.com/aadithyaa9/finance-dashboard/internal/response"
)

type contextKey string

const claimsKey contextKey = "claims"

// Authenticate validates the Bearer JWT and injects claims into context.
func Authenticate(svc *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing or malformed authorization header")
				return
			}

			claims, err := svc.ParseToken(strings.TrimPrefix(header, "Bearer "))
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ClaimsFromContext pulls JWT claims out of the request context.
func ClaimsFromContext(ctx context.Context) *auth.Claims {
	c, _ := ctx.Value(claimsKey).(*auth.Claims)
	return c
}
