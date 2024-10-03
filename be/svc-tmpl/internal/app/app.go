package app

import (
	"context"
	"fmt"
	"io"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

type OptionsFunc func(a *App) error

type App struct {
	blobFetcher blob.Fetcher
}

func New(opts ...OptionsFunc) (*App, error) {
	a := &App{}
	for _, optFunc := range opts {
		err := optFunc(a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func WithBlobFetcher(blobFetcher blob.Fetcher) OptionsFunc {
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

type filePipe struct {
	reader io.ReadCloser
}

func (f *filePipe) WriteTo(w io.Writer) (int64, error) {
	defer f.reader.Close()
	return io.Copy(w, f.reader)
}

func (a *App) GetFile(ctx context.Context, fileId string) (*ddd.AttachmentWriterTo, error) {
	var err error
	ctx, span := logger.Span(ctx, "GetFile", nil)
	defer func() {
		span.End(err)
	}()

	// fileName := "glass-mug-variant.png"
	// p := "/Users/vesko/go/src/github.com/ggrrrr/btmt-ui"
	// file, err := os.Open(fmt.Sprintf("%s/%s", p, fileName))
	// if err != nil {
	// 	return nil, err
	// }

	logger.InfoCtx(ctx).
		Any("info", fileId).
		Msg("Fetch")

	res, err := a.blobFetcher.Fetch(ctx, "localhost", fmt.Sprintf("images/%s", fileId))
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Fetch")
		return nil, err
	}
	logger.InfoCtx(ctx).
		Any("info", res).
		Any("id", fileId).
		Msg("Fetch")

	return &ddd.AttachmentWriterTo{
			ContentType: "image/png",
			Version:     res.Id.Version(),
			Name:        res.Info.Name,
			WriterTo: &filePipe{
				reader: res.ReadCloser,
			},
		},
		nil
}

// func (*App) Render(ctx context.Context, templId string, values any) (ddd.ReanderResponse, error) {
// 	// testReader := strings.NewReader("")

// 	return nil, nil
// }
