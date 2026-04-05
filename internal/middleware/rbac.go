package middleware

import (
	"net/http"

	"github.com/aadithyaa9/finance-dashboard/internal/response"
	"github.com/aadithyaa9/finance-dashboard/internal/users"
)

// RequireRole allows only the listed roles to proceed.
// Must be placed after Authenticate in the middleware chain.
//
//	r.With(RequireRole("admin")).Delete("/records/{id}", h.Delete)
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[users.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[users.Role(r)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil {
				response.Error(w, http.StatusUnauthorized, "unauthenticated")
				return
			}
			if _, ok := allowed[claims.Role]; !ok {
				response.Error(w, http.StatusForbidden, "you do not have permission to perform this action")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
