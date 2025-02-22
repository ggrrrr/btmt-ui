package rest

import (
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	peoplepb "github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb/v1"
)

func (s *server) List(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.List", nil)
	defer func() {
		span.End(err)
	}()

	var req peoplepb.ListRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("List")
		web.SendError(ctx, w, err)
		return
	}

	logger.InfoCtx(ctx).Any("filter", req.String()).Msg("List")
	out, err := s.app.List(ctx, req.ToFilter())
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("List")
		web.SendError(ctx, w, err)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", out)
}

func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.Get", nil)
	defer func() {
		span.End(err)
	}()

	var req peoplepb.GetRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("Get")
		web.SendError(ctx, w, err)
		return

	}
	if req.Id == "" {
		logger.ErrorCtx(ctx, err).Str("error", "empty id").Msg("Get")
		web.SendJSONErrorBadRequest(ctx, w, "empty id", nil)
		return
	}
	logger.InfoCtx(r.Context()).Any("id", &req.Id).Msg("Get")
	p, err := s.app.GetById(ctx, req.Id)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Get")
		web.SendJSONSystemError(ctx, w, "system error, please try again later", err, nil)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", p)
}

func (s *server) Save(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.Save", nil)
	defer func() {
		span.End(err)
	}()

	var req peoplepb.SaveRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Save")
		web.SendError(ctx, w, err)
		return
	}
	logger.InfoCtx(ctx).Any("person", &req).Msg("Save")
	err = s.app.Save(r.Context(), req.Data)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Save")
		web.SendError(ctx, w, err)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", peoplepb.SavePayload{Id: req.Data.Id})
}

func (s *server) Update(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.Update", nil)
	defer func() {
		span.End(err)
	}()

	var req peoplepb.UpdateRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("person", &req).Msg("Update")
		web.SendError(ctx, w, err)
		return
	}

	logger.InfoCtx(ctx).Any("person", &req).Msg("Update")
	err = s.app.Update(r.Context(), req.Data)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("person", &req).Msg("Update")
		web.SendError(ctx, w, err)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", nil)
}
