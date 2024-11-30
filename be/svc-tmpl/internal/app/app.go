package app

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

const (
	tmplBlobFolder string = "templates/images"
	tmplTmplFolder string = "templates/tmpl"
)

type (
	OptionsFunc func(a *App) error

	tmplRepo interface {
		Save(ctx context.Context, template *tmplpb.Template) error
		Update(ctx context.Context, template *tmplpb.Template) error
		List(ctx context.Context, filter app.FilterFactory) (result []*tmplpb.Template, err error)
		GetById(ctx context.Context, fromId string) (*tmplpb.Template, error)
	}

	App struct {
		appPolices   roles.AppPolices
		blobStore    blob.Store
		imagesFolder blob.BlobId
		stateStore   state.StateStore
		repo         tmplRepo
	}
)

func New(opts ...OptionsFunc) (*App, error) {
	imagesFolder, err := blob.ParseBlobDir(tmplBlobFolder)
	if err != nil {
		return nil, err
	}

	a := &App{
		imagesFolder: imagesFolder,
	}
	for _, optFunc := range opts {
		err := optFunc(a)
		if err != nil {
			return nil, err
		}
	}
	if a.appPolices == nil {
		logger.Warn().Msg("use mock AppPolices")
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
	return func(a *App) error {
		a.blobStore = blobStore
		return nil
	}
}

func WithTmplRepo(repo tmplRepo) OptionsFunc {
	return func(a *App) error {
		a.repo = repo
		return nil
	}
}

func WithStateStore(store state.StateStore) OptionsFunc {
	return func(a *App) error {
		a.stateStore = store
		return nil
	}
}
