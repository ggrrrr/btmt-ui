package app

import (
	"bytes"
	"context"
	"fmt"
	htmltemplate "html/template"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/templ"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
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

func (v *tmplValidator) RenderImg(imageName string) htmltemplate.HTML {
	imageId, err := v.app.imagesFolder.SetIdVersionFromString(imageName)
	if err != nil {
		v.errors[imageName] = fmt.Sprintf("name %v", err)
		return htmltemplate.HTML(fmt.Sprintf(`<strong> incorrect image name %s </strong>`, imageName))
	}

	_, err = v.app.blobStore.Head(v.ctx, v.realm, imageId)
	if err != nil {
		v.errors[imageName] = fmt.Sprintf("fetch %v", err)
		return htmltemplate.HTML(fmt.Sprintf(`<strong> image %s not found</strong>`, imageName))
	}

	v.images = append(v.images, imageName)
	suffixUrl := ""

	if v.resized {
		suffixUrl = "/resized"
	}

	return htmltemplate.HTML(fmt.Sprintf(`<img src="http://localhost:8010/tmpl/image/%s%s" ></img>`, imageName, suffixUrl))
}

func (a *App) validate(ctx context.Context, authInfo roles.AuthInfo, template *tmplpb.TemplateUpdate) (*tmplValidator, error) {
	var err error
	ctx, span := logger.Span(ctx, "validate", template)
	defer func() {
		span.End(err)
	}()

	tmplValidator := validator(ctx, authInfo.Realm, a)

	htmlTmpl, err := htmltemplate.New("template_data").
		Funcs(htmltemplate.FuncMap{
			"renderImg": tmplValidator.RenderImg,
		}).
		Parse(template.Body)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = htmlTmpl.Execute(buf, templ.TemplateData{})
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

func (errs TmplError) Error() string {
	if errs == nil {
		return ""
	}
	if len(errs) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for k, v := range errs {
		buffer.WriteString(fmt.Sprintf("file[%s]:%v", k, v))
	}
	return buffer.String()
}
