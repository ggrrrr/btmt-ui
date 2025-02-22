package app

import (
	"bytes"
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/templ"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func (a *Application) RenderHtml(ctx context.Context, render *tmplpb.RenderRequest) (string, error) {
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

	htmlTmpl, err := templ.NewHtml(render.Body, templ.WithRenderImg(tmplValidator.RenderImg))
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("tmpl.parse.body")
		return "", app.BadRequestError("parse body", err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = htmlTmpl.Execute(buf, render.Data)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("tmpl.Execute")
		return "", app.BadRequestError("template execute", err)
	}

	return buf.String(), nil

}
