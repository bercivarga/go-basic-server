package auth

import (
	"encoding/json"
	"net/http"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/middleware"
	"github.com/bercivarga/go-basic-server/internal/router"
	"github.com/bercivarga/go-basic-server/internal/services/user"
	"github.com/bercivarga/go-basic-server/internal/utils"
)

type Handler struct {
	app *app.App
}

func New(a *app.App) *Handler {
	return &Handler{app: a}
}

func (h *Handler) Register(r *router.Router) {
	withAuthMiddleware := router.ComposeMiddleware(middleware.Auth)

	r.HandleFunc(http.MethodPost, "/auth/signup", h.signup)
	r.HandleFunc(http.MethodPost, "/auth/login", h.login)
	r.HandleFunc(http.MethodPost, "/auth/refresh", h.refresh)

	r.HandleFunc(http.MethodPost, "/auth/logout", withAuthMiddleware(h.logout))
}

type SignupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *Handler) signup(a *app.App, w http.ResponseWriter, r *http.Request) {
	var creds SignupRequest
	if err := utils.BindAndValidate(r, &creds); err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	err := a.UserService.CreateUser(r.Context(), user.CreateUserRequest{
		Email:    creds.Email,
		Password: creds.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *Handler) login(a *app.App, w http.ResponseWriter, r *http.Request) {
	var creds LoginRequest
	if err := utils.BindAndValidate(r, &creds); err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	tokens, err := a.AuthService.Login(r.Context(), creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}

type LogoutRequestHeaders struct {
	Authorization string `validate:"required,jwt"`
}

func (h *Handler) logout(a *app.App, w http.ResponseWriter, r *http.Request) {
	token := utils.ExtractBearerToken(r)
	validationData := LogoutRequestHeaders{
		Authorization: token,
	}

	if err := utils.Validate(validationData); err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	err := a.AuthService.Logout(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *Handler) refresh(a *app.App, w http.ResponseWriter, r *http.Request) {
	var body RefreshRequest
	if err := utils.BindAndValidate(r, &body); err != nil {
		utils.RespondWithValidationErrors(w, err)
		return
	}

	tokens, err := a.AuthService.RefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(tokens)
}
