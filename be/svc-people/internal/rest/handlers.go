package rest

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	peoplepb "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

func (s *server) List(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.List")
	defer func() {
		span.End(err)
	}()

	var req peoplepb.ListRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "List")
		web.SendError(ctx, w, err)
		return
	}

	out, err := s.app.List(ctx, req.ToFilter())
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "List")
		web.SendError(ctx, w, err)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", out)
}

func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.Get")
	defer func() {
		span.End(err)
	}()

	var req peoplepb.GetRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		log.Log().ErrorCtx(r.Context(), err, "Get")
		web.SendError(ctx, w, err)
		return

	}
	if req.Id == "" {
		log.Log().ErrorCtx(ctx, fmt.Errorf("empty id"), "Get")
		web.SendJSONErrorBadRequest(ctx, w, "empty id", nil)
		return
	}
	p, err := s.app.GetById(ctx, req.Id)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "Get")
		web.SendJSONSystemError(ctx, w, "system error, please try again later", err, nil)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", p)
}

func (s *server) Save(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.Save")
	defer func() {
		span.End(err)
	}()

	var req peoplepb.SaveRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "Save")
		web.SendError(ctx, w, err)
		return
	}
	log.Log().InfoCtx(ctx, "save", slog.Any("person", &req))
	err = s.app.Save(r.Context(), req.Data)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "Save")
		web.SendError(ctx, w, err)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", peoplepb.SavePayload{Id: req.Data.Id})
}

func (s *server) Update(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.Update")
	defer func() {
		span.End(err)
	}()

	var req peoplepb.UpdateRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "Update", slog.Any("person", &req))
		web.SendError(ctx, w, err)
		return
	}

	log.Log().InfoCtx(ctx, "Update", slog.Any("person", &req))
	err = s.app.Update(r.Context(), req.Data)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "Update", slog.Any("person", &req))
		web.SendError(ctx, w, err)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", nil)
}
