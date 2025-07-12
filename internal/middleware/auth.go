package middleware

import (
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

		claims, err := a.JwtManager.Verify(token)
		if err != nil || !a.SessionStore.IsValid(r.Context(), claims.UserID, token) {
			http.Error(w, "invalid or expired session", http.StatusUnauthorized)
			return
		}

		next(a, w, r)
	}
}
