package user

import (
	"encoding/json"
	"net/http"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/router"
	"github.com/bercivarga/go-basic-server/internal/stores/user"
)

type Handler struct {
	app   *app.App
	store user.UserStore
}

func New(a *app.App) *Handler {
	store := user.NewUserStore(a.DB)
	return &Handler{app: a, store: store}
}

func (h *Handler) Register(r *router.Router) {
	r.HandleFunc("/users", h.list)
}

func (h *Handler) list(a *app.App, w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(users)
}
