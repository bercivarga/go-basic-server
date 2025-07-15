package health

import (
	"fmt"
	"net/http"

	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/router"
)

type Handler struct {
	app *app.App
}

func New(a *app.App) *Handler {
	return &Handler{app: a}
}

func (h *Handler) Register(r *router.Router) {
	r.HandleFunc(http.MethodGet, "/health", h.HealthCheck)
}

func (h *Handler) HealthCheck(a *app.App, w http.ResponseWriter, r *http.Request) {
	a.Logger.Info("Health check")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
