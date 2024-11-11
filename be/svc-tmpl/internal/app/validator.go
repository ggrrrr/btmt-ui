package app

import (
	"context"
	"fmt"
	"html/template"
)

type (
	tmplValidator struct {
		ctx     context.Context
		app     *App
		realm   string
		errors  map[string]error
		resized bool
		images  []string
	}
)

func validator(ctx context.Context, realm string, app *App) *tmplValidator {
	return &tmplValidator{
		ctx:     ctx,
		app:     app,
		realm:   realm,
		errors:  map[string]error{},
		images:  []string{},
		resized: false,
	}
}

func (v *tmplValidator) RenderImg(imageName string) template.HTML {
	imageId, err := v.app.imagesFolder.SetIdVersionFromString(imageName)
	if err != nil {
		v.errors[imageName] = fmt.Errorf("incorrect image[%s] name %w", imageName, err)
		return template.HTML(fmt.Sprintf(`<strong> incorrect image name %s </strong>`, imageName))
	}

	_, err = v.app.blobStore.Head(v.ctx, v.realm, imageId)
	if err != nil {
		v.errors[imageName] = fmt.Errorf("image[%s] not found %w", imageName, err)
		return template.HTML(fmt.Sprintf(`<strong> incorrect image name %s </strong>`, imageName))
	}

	prefix := ""

	if v.resized {
		prefix = "/resized"
	}

	return template.HTML(fmt.Sprintf(`<img src="http://localhost:8010/tmpl/image/%s%s" ></img>`, imageName, prefix))
}
