package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
)

type (
	AppHandler interface {
		Update(w http.ResponseWriter, r *http.Request)
		Save(w http.ResponseWriter, r *http.Request)
		List(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
	}

	server struct {
		app app.App
	}
)

var _ (AppHandler) = (*server)(nil)

func Handler(a app.App) *server {
	return &server{
		app: a,
	}
}

func Router(h AppHandler) chi.Router {
	router := chi.NewRouter()
	router.Post("/update", h.Update)
	router.Post("/save", h.Save)
	router.Post("/list", h.List)
	router.Post("/get", h.Get)

	return router
}
