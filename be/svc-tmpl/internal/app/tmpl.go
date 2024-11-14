package app

import (
	"context"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func (a *App) SaveTmpl(ctx context.Context, tmpl *ddd.Template) error {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "SaveTmpl", tmpl)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return err
	}

	// TODO Verify that the template can be rendered
	// TODO Fetch images and other blobs to verify all is good
	// test if we can parse the body template and data value
	if tmpl.Id == "" {
		tmpl.UpdatedAt = time.Now()
		tmpl.CreatedAt = time.Now()
		err = a.repo.Save(ctx, tmpl)
	} else {
		tmpl.UpdatedAt = time.Now()
		err = a.repo.Update(ctx, tmpl)
	}
	return err
}

func (a *App) ListTmpl(ctx context.Context, filter app.FilterFactory) ([]ddd.Template, error) {
	var err error
	ctx, span := logger.Span(ctx, "ListTmpl", nil)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return nil, err
	}

	logger.InfoCtx(ctx).Msg("ListTmpl")

	result, err := a.repo.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *App) GetTmpl(ctx context.Context, id string) (*ddd.Template, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "GetTmpl", nil, logger.TraceKVString("id", id))
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return nil, err
	}

	result, err := a.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
