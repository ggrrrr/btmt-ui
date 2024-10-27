package rest

import (
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
	"github.com/go-chi/chi/v5"
)

type (
	server struct {
		app app.App
	}
)

func New(a app.App) *server {
	return &server{
		app: a,
	}
}

func (s *server) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/v1/noJson400", noJson400)
	router.Post("/v1/json500", json500)
	router.Get("/v1/noJson400", noJson400)
	router.Get("/v1/json500", json500)

	router.Post("/login/passwd", s.LoginPasswd)
	router.Get("/token/validate", s.TokenValidate)
	router.Post("/token/validate", s.TokenValidate)
	router.Post("/token/refresh", s.TokenRefresh)
	router.Post("/user/list", s.UserList)
	router.Get("/user/list", s.UserList)

	return router
}
