package rest

import (
	"fmt"
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/app"
	"github.com/go-chi/chi"
)

type server struct {
	app *app.App
}

func New(a *app.App) *server {
	return &server{
		app: a,
	}
}

func (s *server) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/render", s.Render)
	router.Get("/attachment/get", s.GetAttachment)

	return router
}

func (s *server) Render(w http.ResponseWriter, r *http.Request) {
}

func (s *server) GetAttachment(w http.ResponseWriter, r *http.Request) {
	attch, err := s.app.GetAttachment(r.Context(), "", "")
	if err != nil {
		web.SendError(r.Context(), w, err)
		return
	}

	w.Header().Add("Content-Type", attch.ContentType)
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", attch.Name))
	//
	w.WriteHeader(http.StatusOK)

	attch.WriterTo.WriteTo(w)

}
