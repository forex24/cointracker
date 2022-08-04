package api

import (
	"net/http"

	"github.com/canhlinh/cointracker/backend"
	chi "github.com/go-chi/chi/v5"
)

type Handler struct {
	App    *backend.App
	Router chi.Router
}

func (h *Handler) ApiHandler(f func(c *Context) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := NewContext(h.App, w, r)
		c.Jsonify(f(c))
	}
}
