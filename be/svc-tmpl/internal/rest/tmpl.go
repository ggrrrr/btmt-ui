package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

type RenderRequest struct {
	Items        map[string]any `json:"items"`
	TemplateBody string         `json:"body"`
}

func (s *server) SaveTmpl(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.SaveTmpl", nil)
	defer func() {
		span.End(err)
	}()
	logger.InfoCtx(r.Context()).Msg("rest.SaveTmpl")

	var template *tmplpb.Template
	err = web.DecodeJsonRequest(r, &template)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	tmplErrors, err := s.app.SaveTmpl(ctx, template)
	if err != nil {
		fmt.Printf("\n\n\t\t%#v \n\n", tmplErrors)
		if tmplErrors != nil {
			web.SendErrorBadRequestWithBody(ctx, w, "validation error", err, tmplErrors)
			return
		}
		web.SendError(ctx, w, err)
		return
	}

	web.SendPayload(ctx, w, "ok", template)
}

func (s *server) ListTmpl(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.ListTmpl", nil)
	logger.InfoCtx(ctx).Msg("rest.ListTmpl")
	defer func() {
		span.End(err)
	}()

	list, err := s.app.ListTmpl(ctx, nil)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	var response = struct {
		List []*tmplpb.Template `json:"list"`
	}{
		List: make([]*tmplpb.Template, 0, len(list)),
	}

	response.List = append(response.List, list...)
	web.SendPayload(ctx, w, "ok", response)
}

func (s *server) GetTmpl(w http.ResponseWriter, r *http.Request) {
	tmplId := chi.URLParam(r, "id")
	var err error
	ctx, span := logger.SpanWithAttributes(r.Context(), "rest.GetTmpl", nil, logger.TraceKVString("id", tmplId))
	defer func() {
		span.End(err)
	}()
	logger.InfoCtx(r.Context()).Msg("rest.GetTmpl")

	tmpl, err := s.app.GetTmpl(ctx, tmplId)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	web.SendPayload(ctx, w, "ok", tmpl)
}

func (s *server) Render(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "Render", nil)
	defer func() {
		span.End(err)
	}()

	request := &tmplpb.RenderRequest{}

	err = web.DecodeJsonRequest(r, request)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	result, err := s.app.RenderHtml(ctx, request)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	out := &tmplpb.RenderResponse{
		Payload: result,
	}

	web.SendPayload(ctx, w, "ok", out)
}
