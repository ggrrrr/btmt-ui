package app

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func (a *App) SaveTmpl(ctx context.Context, tmplUpdate *tmplpb.TemplateUpdate) (string, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "SaveTmpl", tmplUpdate)
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return "", err
	}

	render, err := a.validate(ctx, authInfo, tmplUpdate)
	if err != nil {
		return "", app.BadRequestError("validate", err)
	}
	if len(render.errors) > 0 {
		logger.ErrorCtx(ctx, &render.errors)
		err = fmt.Errorf("validator error(s)")
		return "", render.errors
		// , app.BadRequestError("validate", err)
	}

	tmpl := tmplFromUpdate(tmplUpdate, render)

	if tmpl.Id == "" {
		return a.saveTmpl(ctx, authInfo, tmpl)
	}

	err = a.updateTmpl(ctx, authInfo, tmpl)

	return tmpl.Id, nil
}

func (a *App) ListTmpl(ctx context.Context, filter app.FilterFactory) ([]*tmplpb.Template, error) {
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

func (a *App) GetTmpl(ctx context.Context, id string) (*tmplpb.Template, error) {
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

func (a *App) saveTmpl(ctx context.Context, authInfo roles.AuthInfo, tmpl *tmplpb.Template) (string, error) {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "saveTmpl", tmpl)
	defer func() {
		span.End(err)
	}()

	err = a.repo.Save(ctx, tmpl)
	if err != nil {
		return "", err
	}

	err = a.uploadTmplBody(ctx, authInfo, tmpl)
	if err != nil {
		return "", err
	}

	return tmpl.Id, nil
}

func (a *App) updateTmpl(ctx context.Context, authInfo roles.AuthInfo, tmpl *tmplpb.Template) error {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "updateTmpl", tmpl)
	defer func() {
		span.End(err)
	}()

	updateBlob := false

	oldTmpl, err := a.repo.GetById(ctx, tmpl.Id)
	if err != nil {
		return err
	}

	if oldTmpl.Body != tmpl.Body {
		updateBlob = true
	}

	err = a.repo.Update(ctx, tmpl)
	if err != nil {
		return err
	}

	if updateBlob {
		err = a.uploadTmplBody(ctx, authInfo, tmpl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) uploadTmplBody(ctx context.Context, authInfo roles.AuthInfo, tmpl *tmplpb.Template) error {
	var err error

	value, err := proto.Marshal(tmpl)
	if err != nil {
		return app.SystemError("cant marshal template", err)
	}

	_, err = a.stateStore.Push(ctx, state.NewEntity{
		Key:   tmpl.Id,
		Value: value,
	})
	if err != nil {
		return app.SystemError("cant push template to blob store", err)
	}

	return nil
}

func tmplFromUpdate(tmplUpdate *tmplpb.TemplateUpdate, validator *tmplValidator) *tmplpb.Template {
	tmtpl := &tmplpb.Template{
		ContentType: tmplUpdate.ContentType,
		Name:        tmplUpdate.Name,
		Labels:      tmplUpdate.Labels,
		Body:        tmplUpdate.Body,
		Images:      validator.images,
	}

	if tmplUpdate.Id != "" {
		tmtpl.Id = tmplUpdate.Id
	}

	return tmtpl
}
