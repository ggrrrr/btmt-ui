package app

import (
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	OptionsFunc func(a *App) error

	App struct {
		appPolices   roles.AppPolices
		blobStore    blob.Store
		imagesFolder blob.BlobId
	}
)

func New(opts ...OptionsFunc) (*App, error) {

	imagesFolder, err := blob.ParseBlobId("images")
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
	return a, nil
}

func WithBlobStore(blobStore blob.Store) OptionsFunc {
	return func(a *App) error {
		a.blobStore = blobStore
		return nil
	}
}
