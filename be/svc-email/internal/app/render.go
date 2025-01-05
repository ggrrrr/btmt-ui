package app

import (
	"context"
	"fmt"
	htmltemplate "html/template"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
)

type (
	templRender struct {
		ctx    context.Context
		app    *Application
		realm  string
		errors map[string]error
		images map[string]blob.BlobReader
	}
)

func createImageRender(ctx context.Context, realm string, app *Application) *templRender {
	return &templRender{
		ctx:    ctx,
		app:    app,
		realm:  realm,
		errors: map[string]error{},
		images: map[string]blob.BlobReader{},
	}
}

func (templ *templRender) renderImg(imageName string) htmltemplate.HTML {
	imageId, err := templ.app.imagesFolder.SetIdVersionFromString(imageName)
	if err != nil {
		templ.errors[imageName] = fmt.Errorf("image id %w", err)
		return htmltemplate.HTML(fmt.Sprintf(`<strong> incorrect image name %s </strong>`, imageName))
	}

	reader, err := templ.app.blobStore.Fetch(templ.ctx, templ.realm, imageId)
	if err != nil {
		templ.errors[imageName] = fmt.Errorf("fetch %w", err)
		return htmltemplate.HTML(fmt.Sprintf(`<strong> fetch image name %s </strong>`, imageName))
	}
	templ.images[imageName] = reader
	// <img src="cid:glass-mug-variant.png" alt="My image" />
	return htmltemplate.HTML(fmt.Sprintf(`<img src="cid:%s" alt="%s"></img>`, imageName, imageName))
}
