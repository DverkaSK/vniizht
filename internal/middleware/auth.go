package middleware

import (
	"context"
	"net/http"
	"strings"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/service"
)

type contextKey string

const ctxUser contextKey = "user"

func Authenticate(auth *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token != "" {
				if user, err := auth.Authenticate(r.Context(), token); err == nil {
					r = r.WithContext(context.WithValue(r.Context(), ctxUser, user))
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if CurrentUser(r) == nil {
			errs.Write(w, http.StatusUnauthorized, errs.Unauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRole(roles ...model.Role) func(http.Handler) http.Handler {
	allowed := make(map[model.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := CurrentUser(r)
			if u == nil {
				errs.Write(w, http.StatusUnauthorized, errs.Unauthorized)
				return
			}
			if _, ok := allowed[u.Role]; !ok {
				errs.Write(w, http.StatusForbidden, errs.Forbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CurrentUser(r *http.Request) *model.User {
	u, _ := r.Context().Value(ctxUser).(*model.User)
	return u
}

func extractToken(r *http.Request) string {
	if c, err := r.Cookie("session"); err == nil {
		return c.Value
	}
	if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return ""
}
