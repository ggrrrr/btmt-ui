package rest

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/common/web"
	tmplpb "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

type RenderRequest struct {
	Items        map[string]any `json:"items"`
	TemplateBody string         `json:"body"`
}

func (s *server) SaveTmpl(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.SaveTmpl")
	defer func() {
		span.End(err)
	}()

	var template *tmplpb.TemplateUpdate
	err = web.DecodeJsonRequest(r, &template)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	id, err := s.app.SaveTmpl(ctx, template)
	if err != nil {
		// TODO work on form fields error
		// fmt.Printf("\n\n\t\t%#v \n\n", tmplErrors)
		// if tmplErrors != nil {
		// 	web.SendErrorBadRequestWithBody(ctx, w, "validation error", err, tmplErrors)
		// 	return
		// }
		web.SendError(ctx, w, err)
		return
	}
	template.Id = id
	web.SendJSONPayload(ctx, w, "ok", template)
}

func (s *server) ListTmpl(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.ListTmpl")
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
	web.SendJSONPayload(ctx, w, "ok", response)
}

func (s *server) GetTmpl(w http.ResponseWriter, r *http.Request) {
	tmplId := chi.URLParam(r, "id")
	var err error
	ctx, span := s.tracer.SpanWithAttributes(r.Context(), "rest.GetTmpl", slog.String("id", tmplId))
	defer func() {
		span.End(err)
	}()

	tmpl, err := s.app.GetTmpl(ctx, tmplId)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}

	web.SendJSONPayload(ctx, w, "ok", tmpl)
}

func (s *server) Render(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "Render")
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

	web.SendJSONPayload(ctx, w, "ok", out)
}
