package rest

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

type (
	tmplApp interface {
		SaveTmpl(ctx context.Context, tmplUpdate *tmplpb.TemplateUpdate) (string, error)
		ListTmpl(ctx context.Context, filter app.FilterFactory) ([]*tmplpb.Template, error)
		GetTmpl(ctx context.Context, id string) (*tmplpb.Template, error)
		RenderHtml(ctx context.Context, render *tmplpb.RenderRequest) (string, error)

		// images
		GetImage(ctx context.Context, fileId string, maxWight int) (*ddd.FileWriterTo, error)
		PutImage(ctx context.Context, tempFile blob.TempFile) error
		ListImages(ctx context.Context) ([]ddd.ImageInfo, error)
		GetResizedImage(ctx context.Context, fileId string) (*ddd.FileWriterTo, error)
	}
	server struct {
		app tmplApp
	}
)

func New(a tmplApp) *server {
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
