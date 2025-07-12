package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/auth"
	"github.com/bercivarga/go-basic-server/internal/middleware"
	"github.com/bercivarga/go-basic-server/internal/router"
	"github.com/bercivarga/go-basic-server/internal/stores/user"
	"github.com/bercivarga/go-basic-server/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	JWT_DURATION = 24 * time.Hour * 7 // 7 days
)

type Handler struct {
	app        *app.App
	store      *user.Store
	jwtManager *auth.JWTManager
}

func New(a *app.App) *Handler {
	store := user.NewStore(a.DB)

	jwtManager := auth.NewJWTManager(a.Config.JWTSecret, JWT_DURATION)

	return &Handler{app: a, store: store, jwtManager: jwtManager}
}

func (h *Handler) Register(r *router.Router) {
	withAuthMiddleware := utils.ComposeMiddleware(middleware.Auth)

	r.HandleFunc("/login", h.login)
	r.HandleFunc("/signup", h.signup)
	r.HandleFunc("/me", withAuthMiddleware(h.me))
	r.HandleFunc("/users", withAuthMiddleware(h.list))
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

	u, err := h.store.GetByEmail(r.Context(), creds.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if !utils.CheckPasswordHash(creds.Password, u.PasswordHash) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	token, err := h.jwtManager.Generate(u.ID)
	if err != nil {
		http.Error(w, "could not create token", http.StatusInternalServerError)
		return
	}

	expires := time.Now().Add(h.jwtManager.TokenDuration)
	_ = h.app.SessionStore.Create(r.Context(), u.ID, token, expires)

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *Handler) me(a *app.App, w http.ResponseWriter, r *http.Request) {
	// Example protected route
	w.Write([]byte("You are authenticated"))
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
