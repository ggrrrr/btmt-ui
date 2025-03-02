package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.be.svc-auth"

type (
	AppHandler interface {
		LoginPasswd(w http.ResponseWriter, r *http.Request)
		TokenValidate(w http.ResponseWriter, r *http.Request)
		TokenRefresh(w http.ResponseWriter, r *http.Request)
		UserList(w http.ResponseWriter, r *http.Request)
	}
	server struct {
		tracer tracer.OTelTracer
		app    app.App
	}
)

func Handler(a app.App) *server {
	return &server{
		tracer: tracer.Tracer(otelScope),
		app:    a,
	}
}

func Router(h AppHandler) chi.Router {
	router := chi.NewRouter()

	router.Post("/login/passwd", h.LoginPasswd)

	router.Get("/token/validate", h.TokenValidate)
	router.Post("/token/validate", h.TokenValidate)

	router.Post("/token/refresh", h.TokenRefresh)

	router.Post("/login/list", h.UserList)
	router.Get("/login/list", h.UserList)

	return router
}
