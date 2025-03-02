package app

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.be.svc-tmpl"

const (
	tmplBlobFolder string = "templates/images"
	tmplTmplFolder string = "templates/tmpl"
)

type (
	OptionsFunc func(a *Application) error

	tmplRepo interface {
		Save(ctx context.Context, template *tmplpb.Template) error
		Update(ctx context.Context, template *tmplpb.Template) error
		List(ctx context.Context, filter app.FilterFactory) (result []*tmplpb.Template, err error)
		GetById(ctx context.Context, fromId string) (*tmplpb.Template, error)
	}

	Application struct {
		tracer       tracer.OTelTracer
		appPolices   roles.AppPolices
		blobStore    blob.Store
		imagesFolder blob.BlobId
		stateStore   state.StateStore
		repo         tmplRepo
	}

	App interface {
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
)

var _ (App) = (*Application)(nil)

func New(opts ...OptionsFunc) (*Application, error) {
	imagesFolder, err := blob.ParseBlobDir(tmplBlobFolder)
	if err != nil {
		return nil, err
	}

	a := &Application{
		tracer:       tracer.Tracer(otelScope),
		imagesFolder: imagesFolder,
	}
	for _, optFunc := range opts {
		err := optFunc(a)
		if err != nil {
			return nil, err
		}
	}
	if a.appPolices == nil {
		log.Log().Warn(nil, "use mock AppPolices")
		a.appPolices = roles.NewAppPolices()
	}
	if a.blobStore == nil {
		return nil, fmt.Errorf("blobStore is nil")
	}

	if a.appPolices == nil {
		return nil, fmt.Errorf("appPolices is nil")
	}

	if a.blobStore == nil {
		return nil, fmt.Errorf("blobStore is nil")
	}

	if a.repo == nil {
		return nil, fmt.Errorf("repo is nil")
	}

	if a.stateStore == nil {
		return nil, fmt.Errorf("stateStore is nil")
	}

	return a, nil
}

func WithBlobStore(blobStore blob.Store) OptionsFunc {
	return func(a *Application) error {
		a.blobStore = blobStore
		return nil
	}
}

func WithTmplRepo(repo tmplRepo) OptionsFunc {
	return func(a *Application) error {
		a.repo = repo
		return nil
	}
}

func WithStateStore(store state.StateStore) OptionsFunc {
	return func(a *Application) error {
		a.stateStore = store
		return nil
	}
}
