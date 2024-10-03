package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-events/internal/app"
)

type (
	server struct {
		app app.App
	}
)

// func New() *server {
func New(a app.App) *server {
	return &server{
		app: a,
	}
}

func (s *server) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/v1/tmpl", s.Tmpl)
	return router
}

func (s *server) Tmpl(w http.ResponseWriter, r *http.Request) {
	web.SendPayload(r.Context(), w, "ok", "tmpl")
}
