package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-tmpl/internal/ddd"
)

func (a *App) SaveTmpl(ctx context.Context, tmpl *ddd.Template) error {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "SaveTmpl", tmpl)
	defer func() {
		span.End(err)
	}()

	// Fetch images and other blobs to verify all is good
	// test if we can parse the body with template
	if tmpl.Id == "" {
		err = a.repo.Save(ctx, tmpl)
	} else {
		err = a.repo.Update(ctx, tmpl)
	}
	if err != nil {

	}
	return err

}

func (a *App) ListTmpl(ctx context.Context, filter app.FilterFactory) ([]ddd.Template, error) {
	var err error
	ctx, span := logger.Span(ctx, "ListTmpl", nil)
	defer func() {
		span.End(err)
	}()

	logger.InfoCtx(ctx).Msg("ListTmpl")

	result, err := a.repo.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *App) GetTmpl(ctx context.Context, id string) (*ddd.Template, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "GetTmpl", nil, logger.KVString("id", id))
	defer func() {
		span.End(err)
	}()

	result, err := a.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
