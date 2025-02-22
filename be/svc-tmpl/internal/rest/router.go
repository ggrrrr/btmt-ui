package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/app"
)

type (
	appHandler interface {
		SaveTmpl(w http.ResponseWriter, r *http.Request)
		ListTmpl(w http.ResponseWriter, r *http.Request)
		GetTmpl(w http.ResponseWriter, r *http.Request)
		Render(w http.ResponseWriter, r *http.Request)
		// images
		UploadImage(w http.ResponseWriter, r *http.Request)
		GetImage(w http.ResponseWriter, r *http.Request)
		GetResizedImage(w http.ResponseWriter, r *http.Request)
		ListImages(w http.ResponseWriter, r *http.Request)
	}

	server struct {
		app app.App
	}
)

func Handler(a app.App) *server {
	return &server{
		app: a,
	}
}

func Router(s appHandler) chi.Router {
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
