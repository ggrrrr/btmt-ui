package rest

import (
	"fmt"
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
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
	router.Get("/file/{id}", s.GetFile)
	router.Get("/tmpl/get", s.GetTmpl)

	return router
}

func (s *server) GetTmpl(w http.ResponseWriter, r *http.Request) {
	logger.InfoCtx(r.Context()).Msg("rest.GetTmpl")

	web.SendError(r.Context(), w, fmt.Errorf("asdasd"))
}

func (s *server) Render(w http.ResponseWriter, r *http.Request) {
}

func (s *server) GetFile(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.GetFile", nil)
	defer func() {
		span.End(err)
	}()

	fileId := chi.URLParam(r, "id")
	download := r.URL.Query().Get("download")
	logger.DebugCtx(ctx).
		Str("id", fileId).
		Msg("rest.GetFile")

	attch, err := s.app.GetFile(ctx, fileId)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	w.Header().Add("Content-Type", attch.ContentType)

	if download != "" {
		w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", attch.Name))
	}

	w.WriteHeader(http.StatusOK)

	attch.WriterTo.WriteTo(w)

}
