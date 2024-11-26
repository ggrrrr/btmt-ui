package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
)

type (
	server struct {
		app *app.App
	}
)

func New(a *app.App) *server {
	return &server{
		app: a,
	}
}

func (s *server) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/update", s.Update)
	router.Post("/save", s.Save)
	router.Post("/list", s.List)
	router.Post("/get", s.Get)

	return router
}
