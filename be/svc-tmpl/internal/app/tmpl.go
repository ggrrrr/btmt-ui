package app

import (
	"context"
	"strings"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	appError "github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/blob"
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

	if err != nil {
		return err
	}

	err = a.saveTmpl(ctx, authInfo, tmpl)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) saveTmpl(ctx context.Context, authInfo roles.AuthInfo, tmpl *ddd.Template) error {
	tmplBlobId, err := a.tmplFolder.SetIdVersionFromString(tmpl.Id)
	if err != nil {
		return appError.SystemError("cant create template blob name", err)
	}

	buffer := strings.NewReader(tmpl.Body)

	_, err = a.blobStore.Push(ctx, authInfo.Realm, tmplBlobId, blob.BlobMD{
		Type:          "template",
		ContentType:   tmpl.ContentType,
		Name:          tmpl.Name,
		ContentLength: int64(len(tmpl.Body)),
	}, buffer)

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
