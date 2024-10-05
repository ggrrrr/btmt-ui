package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

type (
	OptionsFunc func(a *App) error

	App struct {
		appPolices   roles.AppPolices
		blobFetcher  blob.Store
		imagesFolder string
	}
)

func New(opts ...OptionsFunc) (*App, error) {
	a := &App{
		imagesFolder: "images",
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
	return a, nil
}

func WithBlobStore(blobFetcher blob.Store) OptionsFunc {
	return func(a *App) error {
		a.blobFetcher = blobFetcher
		return nil
	}
}

func (a *App) GetTmpl(ctx context.Context, tmplId string, tmplVersion string) (*ddd.Tmpl, error) {
	var err error
	ctx, span := logger.Span(ctx, "GetTmpl", nil)
	defer func() {
		span.End(err)
	}()

	logger.InfoCtx(ctx).Str("tmeplId", tmplId).Msg("GetTmpl")

	result, err := a.blobFetcher.Head(ctx, "localhost", "images/beer1")
	if err != nil {
		return nil, err
	}
	logger.InfoCtx(ctx).Any("info", result).Msg("got")
	return &ddd.Tmpl{
		ContentType: "text/markdown",
		Body: `# Header1 {{ .UserInfo.User }}
		## item 1 {{  }}
		## Footers`,
	}, nil
}

// func (*App) Render(ctx context.Context, templId string, values any) (ddd.ReanderResponse, error) {
// 	// testReader := strings.NewReader("")

// 	return nil, nil
// }
