package app

import (
	"bytes"
	"context"
	htmltemplate "html/template"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func (a *App) RenderHtml(ctx context.Context, render *tmplpb.RenderRequest) (string, error) {
	var err error
	ctx, span := logger.Span(ctx, "RenderHtml", nil)
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
		Parse(render.Body)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("tmpl.parse.body")
		return "", app.BadRequestError("parse body", err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buf, render.Data)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("tmpl.Execute")
		return "", app.BadRequestError("template execute", err)
	}

	return buf.String(), nil

}
