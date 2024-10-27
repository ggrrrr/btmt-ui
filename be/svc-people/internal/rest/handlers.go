package rest

import (
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb"
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
	people := []*peoplepb.Person{}
	for _, p := range out {
		person := peoplepb.FromPerson(&p)
		logger.DebugCtx(r.Context()).Any("person from", p).Msg("ListResult")
		logger.DebugCtx(r.Context()).Any("person to", person).Msg("ListResult")
		people = append(people, person)
	}
	web.SendPayload(ctx, w, "ok", people)
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
		web.SendErrorBadRequest(ctx, w, "empty id", nil)
		return
	}
	logger.InfoCtx(r.Context()).Any("id", &req.Id).Msg("Get")
	p, err := s.app.GetById(ctx, req.Id)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Get")
		web.SendSystemError(ctx, w, "system error, please try again later", err, nil)
		return
	}
	web.SendPayload(ctx, w, "ok", p)
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
	p := req.ToPerson()
	logger.InfoCtx(ctx).Any("person", &req).Msg("Save")
	err = s.app.Save(r.Context(), p)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("Save")
		web.SendError(ctx, w, err)
		return
	}
	web.SendPayload(ctx, w, "ok", peoplepb.SavePayload{Id: p.Id})
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
	p := req.ToPerson()
	logger.InfoCtx(ctx).Any("person", &req).Msg("Update")
	err = s.app.Update(r.Context(), p)
	if err != nil {
		logger.ErrorCtx(ctx, err).Any("person", &req).Msg("Update")
		web.SendError(ctx, w, err)
		return
	}
	web.SendPayload(ctx, w, "ok", nil)
}
