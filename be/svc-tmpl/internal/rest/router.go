package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/app"
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
	router.Get("/image/{id}", s.GetImage)
	// router.Get("/image/some/*", s.GetImage)
	// router.Get("/image/resized/*", s.GetResizedImage)
	router.Get("/image/{id}/resized", s.GetResizedImage)
	router.Post("/image", s.UploadImage)
	router.Get("/images", s.ListImages)

	router.Post("/manage/render", s.Render)
	router.Post("/manage/list", s.ListTmpl)
	router.Get("/manage/{id}", s.GetTmpl)
	router.Post("/manage/save", s.SaveTmpl)

	return router
}
