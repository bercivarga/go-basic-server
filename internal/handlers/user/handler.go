package user

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/middleware"
	"github.com/bercivarga/go-basic-server/internal/router"
	"github.com/bercivarga/go-basic-server/internal/stores/user"
	"github.com/bercivarga/go-basic-server/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	app   *app.App
	store *user.Store
}

func New(a *app.App) *Handler {
	store := user.NewStore(a.DB)

	return &Handler{app: a, store: store}
}

func (h *Handler) Register(r *router.Router) {
	withAuthMiddleware := router.ComposeMiddleware(middleware.Auth)
	withAdminMiddleware := router.ComposeMiddleware(
		middleware.Auth,
		middleware.AdminOnly,
	)

	r.HandleFunc("/login", h.login)
	r.HandleFunc("/signup", h.signup)
	r.HandleFunc("/refresh", h.refresh)
	r.HandleFunc("/logout", h.logout)

	r.HandleFunc("/me", withAuthMiddleware(h.me))
	r.HandleFunc("/users", withAdminMiddleware(h.list))
}

func (h *Handler) signup(a *app.App, w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "hashing failed", http.StatusInternalServerError)
		return
	}

	_, err = h.store.Create(r.Context(), creds.Email, string(hash))
	if err != nil {
		http.Error(w, "user exists or db error", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) login(a *app.App, w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	tokens, err := a.AuthService.Login(r.Context(), creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

func (h *Handler) logout(a *app.App, w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	err := a.AuthService.SessionStore.DeleteByToken(r.Context(), token)
	if err != nil {
		http.Error(w, "could not delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) refresh(a *app.App, w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.RefreshToken == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Validate refresh token
	session, err := a.AuthService.SessionStore.GetByRefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Create new token pair
	accessToken, err := a.AuthService.JwtManager.Generate(session.UserID)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "refresh token generation failed", http.StatusInternalServerError)
		return
	}

	// Delete old session and insert new
	if err := a.AuthService.SessionStore.DeleteByRefreshToken(r.Context(), body.RefreshToken); err != nil {
		http.Error(w, "session cleanup failed", http.StatusInternalServerError)
		return
	}

	accessTokenExpireAt, refreshTokenExpireAt := a.AuthService.JwtManager.CreateExpiry()

	err = a.AuthService.SessionStore.Create(r.Context(), session.ID, accessToken, newRefreshToken, accessTokenExpireAt, refreshTokenExpireAt)
	if err != nil {
		http.Error(w, "session creation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *Handler) me(a *app.App, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserIdFromContext(ctx)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	user, err := h.store.GetByID(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	meData := struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}{
		Email: user.Email,
		Role:  user.Role,
	}

	err = json.NewEncoder(w).Encode(meData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func (h *Handler) list(a *app.App, w http.ResponseWriter, r *http.Request) {
	var limit, offset int64 = 10, 0 // TODO: Implement pagination
	users, err := h.store.GetAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
