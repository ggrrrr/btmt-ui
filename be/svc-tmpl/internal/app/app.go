package app

import (
	"context"
	"fmt"
	"os"

	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

type App struct {
}

func New() (*App, error) {
	return &App{}, nil
}

func (a *App) GetTmpl(ctx context.Context, tmpllId string, tmplVersion string) (ddd.Tmpl, error) {
	return ddd.Tmpl{
		ContentType: "text/markdown",
		Body: `# Header1 {{ .UserInfo.User }}
		## item 1 {{  }}
		## Footers`,
	}, nil
}

func (a *App) GetAttachment(ctx context.Context, tmpllId string, tmplVersion string) (*ddd.AttachmentWriterTo, error) {
	fileName := "glass-mug-variant.png"
	p := "/Users/vesko/go/src/github.com/ggrrrr/btmt-ui"
	file, err := os.Open(fmt.Sprintf("%s/%s", p, fileName))
	if err != nil {
		return nil, err
	}

	return &ddd.AttachmentWriterTo{
			ContentType: "image/png",
			Version:     "v1",
			Name:        "glass-mug-variant.png",
			WriterTo:    file,
		},
		nil
}

// func (*App) Render(ctx context.Context, templId string, values any) (ddd.ReanderResponse, error) {
// 	// testReader := strings.NewReader("")

// 	return nil, nil
// }
