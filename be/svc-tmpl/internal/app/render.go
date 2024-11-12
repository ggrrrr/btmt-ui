package app

import (
	"bytes"
	"context"
	htmltemplate "html/template"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func (a *App) RenderHtml(ctx context.Context, template ddd.Template, data ddd.TemplateData) (string, error) {
	var err error
	ctx, span := logger.Span(ctx, "RenderHtml", template)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return "", err
	}

	tmplValidator := validator(ctx, authInfo.Realm, a)
	tmplValidator.resized = true

	tmpl, err := htmltemplate.New("template_data").
		Funcs(htmltemplate.FuncMap{
			"renderImg": tmplValidator.RenderImg,
		}).
		Parse(template.Body)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, data)

	if len(tmplValidator.errors) > 0 {
		logger.ErrorCtx(ctx, tmplValidator.Errros()[0]).Errs("errors", tmplValidator.Errros()).Msg("validator")
	}

	return buf.String(), err

}
