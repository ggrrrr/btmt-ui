package app

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func (a *Application) SaveTmpl(ctx context.Context, tmplUpdate *tmplpb.TemplateUpdate) (string, error) {
	var err error
	ctx, span := a.tracer.SpanWithData(ctx, "SaveTmpl", tmplUpdate)
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

func (a *Application) ListTmpl(ctx context.Context, filter app.FilterFactory) ([]*tmplpb.Template, error) {
	var err error
	ctx, span := a.tracer.Span(ctx, "ListTmpl")
	defer func() {
		span.End(err)
	}()

	authInfo := roles.AuthInfoFromCtx(ctx)
	err = a.appPolices.CanDo(authInfo.Realm, "some", authInfo)
	if err != nil {
		return nil, err
	}

	result, err := a.repo.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *Application) GetTmpl(ctx context.Context, id string) (*tmplpb.Template, error) {
	var err error
	ctx, span := a.tracer.SpanWithAttributes(ctx, "GetTmpl", slog.String("id", id))
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

func (a *Application) saveTmpl(ctx context.Context, authInfo roles.AuthInfo, tmpl *tmplpb.Template) (string, error) {
	var err error
	ctx, span := a.tracer.SpanWithData(ctx, "saveTmpl", tmpl)
	defer func() {
		span.End(err)
	}()

	err = a.repo.Save(ctx, tmpl)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "repo.Save")
		return "", err
	}

	err = a.uploadTmplBody(ctx, authInfo, tmpl)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "uploadTmplBody")
		return "", err
	}

	return tmpl.Id, nil
}

func (a *Application) updateTmpl(ctx context.Context, authInfo roles.AuthInfo, tmpl *tmplpb.Template) error {
	var err error
	ctx, span := a.tracer.SpanWithData(ctx, "updateTmpl", tmpl)
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

// TODO use authInfo.realm to push into state store
func (a *Application) uploadTmplBody(ctx context.Context, authInfo roles.AuthInfo, tmpl *tmplpb.Template) error {
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
