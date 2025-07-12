package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/router"
)

func Auth(next router.HandleFuncWithApp) router.HandleFuncWithApp {
	return func(a *app.App, w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := a.AuthService.JwtManager.Verify(token)
		if err != nil || !a.AuthService.SessionStore.IsValid(r.Context(), claims.UserID, token) {
			http.Error(w, "invalid or expired session", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey, claims.UserID)

		next(a, w, r.WithContext(ctx))
	}
}

func AdminOnly(next router.HandleFuncWithApp) router.HandleFuncWithApp {
	return func(a *app.App, w http.ResponseWriter, r *http.Request) {
		userId, ok := GetUserIdFromContext(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		err := a.AuthService.CheckRole(r.Context(), userId, "admin")
		if err != nil {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next(a, w, r)
	}
}
