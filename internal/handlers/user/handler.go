package user

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/middleware"
	"github.com/bercivarga/go-basic-server/internal/router"
	"github.com/bercivarga/go-basic-server/internal/services/user"
)

type Handler struct {
	app *app.App
}

func New(a *app.App) *Handler {
	return &Handler{app: a}
}

func (h *Handler) Register(r *router.Router) {
	withAuthMiddleware := router.ComposeMiddleware(middleware.Auth)
	withAdminMiddleware := router.ComposeMiddleware(
		middleware.Auth,
		middleware.AdminOnly,
	)

	r.HandleFunc("/users/me", withAuthMiddleware(h.me))
	r.HandleFunc("/users/list", withAdminMiddleware(h.list))
}

func (h *Handler) me(a *app.App, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := middleware.GetUserIdFromContext(ctx)
	if !ok {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	user, err := a.UserService.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) list(a *app.App, w http.ResponseWriter, r *http.Request) {
	var limit, offset int64

	limit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		limit = 10
	}

	offset, err = strconv.ParseInt(r.URL.Query().Get("offset"), 10, 64)
	if err != nil {
		offset = 0
	}

	users, err := a.UserService.ListUsers(r.Context(), user.ListUsersRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
