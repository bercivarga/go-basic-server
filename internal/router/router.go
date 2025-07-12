package router

import (
	"net/http"

	"github.com/bercivarga/go-basic-server/internal/app"
)

type Router struct {
	app *app.App
	mux *http.ServeMux
}

// New returns a Router ready to plug into http.Server.
func New(a *app.App) *Router {
	r := &Router{
		app: a,
		mux: http.NewServeMux(),
	}
	return r
}

// ServeHTTP lets Router satisfy http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// Handle registers a path ➜ handler pair.
func (r *Router) Handle(pattern string, h http.Handler) {
	r.mux.Handle(pattern, h)
}

// HandleFuncWithApp describes a handler that captures *app.App.
type HandleFuncWithApp func(*app.App, http.ResponseWriter, *http.Request)

// HandleFunc shortcuts to http.HandlerFunc and captures *app.App.
func (r *Router) HandleFunc(pattern string, fn HandleFuncWithApp) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		fn(r.app, w, req)
	})
}
