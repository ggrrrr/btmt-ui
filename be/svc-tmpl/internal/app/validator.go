package app

import (
	"bytes"
	"context"
	"fmt"
	htmltemplate "html/template"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

type (
	tmplValidator struct {
		ctx     context.Context
		app     *App
		realm   string
		errors  TmplError
		resized bool
		images  []string
	}

	TmplError map[string]string
)

func (a *App) validate(ctx context.Context, authInfo roles.AuthInfo, template *ddd.Template) (*tmplValidator, error) {
	var err error
	ctx, span := logger.Span(ctx, "validate", template)
	defer func() {
		span.End(err)
	}()

	tmplValidator := validator(ctx, authInfo.Realm, a)

	tmpl, err := htmltemplate.New("template_data").
		Funcs(htmltemplate.FuncMap{
			"renderImg": tmplValidator.RenderImg,
		}).
		Parse(template.Body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, ddd.TemplateData{})
	if err != nil {
		return nil, err
	}
	return tmplValidator, nil

}

func validator(ctx context.Context, realm string, app *App) *tmplValidator {
	return &tmplValidator{
		ctx:     ctx,
		app:     app,
		realm:   realm,
		errors:  map[string]string{},
		images:  []string{},
		resized: false,
	}
}

func (v *tmplValidator) RenderImg(imageName string) htmltemplate.HTML {
	imageId, err := v.app.imagesFolder.SetIdVersionFromString(imageName)
	if err != nil {
		v.errors[imageName] = fmt.Sprintf("image name:[%s] %v", imageName, err)
		return htmltemplate.HTML(fmt.Sprintf(`<strong> incorrect image name %s </strong>`, imageName))
	}

	_, err = v.app.blobStore.Head(v.ctx, v.realm, imageId)
	if err != nil {
		v.errors[imageName] = fmt.Sprintf("fetch image:[%s]  %v", imageName, err)
		return htmltemplate.HTML(fmt.Sprintf(`<strong> incorrect image name %s </strong>`, imageName))
	}

	suffixUrl := ""

	if v.resized {
		suffixUrl = "/resized"
	}

	return htmltemplate.HTML(fmt.Sprintf(`<img src="http://localhost:8010/tmpl/image/%s%s" ></img>`, imageName, suffixUrl))
}
