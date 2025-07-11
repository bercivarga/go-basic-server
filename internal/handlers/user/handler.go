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
	store *user.Store
}

func New(a *app.App) *Handler {
	store := user.NewStore(a.DB)
	return &Handler{app: a, store: store}
}

func (h *Handler) Register(r *router.Router) {
	r.HandleFunc("/users", h.list)
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
