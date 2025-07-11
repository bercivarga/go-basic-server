package wire

import (
	"github.com/bercivarga/go-basic-server/internal/app"
	"github.com/bercivarga/go-basic-server/internal/handlers/health"
	"github.com/bercivarga/go-basic-server/internal/handlers/user"
	"github.com/bercivarga/go-basic-server/internal/router"
)

// Components collects every feature-handler the service owns.
type Components struct {
	User   *user.Handler
	Health *health.Handler
}

// New builds all handlers that need *app.App.
func New(a *app.App) *Components {
	return &Components{
		User:   user.New(a),
		Health: health.New(a),
	}
}

// RegisterRoutes attaches every handler’s routes to the router.
func (c *Components) RegisterRoutes(r *router.Router) {
	c.User.Register(r)
	c.Health.Register(r)
}
