package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.be.svc-people"

type (
	AppHandler interface {
		Update(w http.ResponseWriter, r *http.Request)
		Save(w http.ResponseWriter, r *http.Request)
		List(w http.ResponseWriter, r *http.Request)
		Get(w http.ResponseWriter, r *http.Request)
	}

	server struct {
		tracer tracer.OTelTracer
		app    app.App
	}
)

var _ (AppHandler) = (*server)(nil)

func Handler(a app.App) *server {
	return &server{
		tracer: tracer.Tracer(otelScope),
		app:    a,
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
