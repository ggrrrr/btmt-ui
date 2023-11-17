package rest

import (
	"net/http"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/peoplepb"
	"github.com/go-chi/chi"
)

type (
	server struct {
		app app.App
	}
)

func New(a app.App) *server {
	return &server{
		app: a,
	}
}

func (s *server) Router() chi.Router {
	router := chi.NewRouter()
	router.Post("/v1/people/update", s.Update)
	router.Post("/v1/people/save", s.Save)
	router.Post("/v1/people/list", s.List)
	router.Post("/v1/people/get", s.Get)

	return router
}

func (s *server) List(w http.ResponseWriter, r *http.Request) {
	var req peoplepb.ListRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(w, err)
	}
	logger.Log().Info().Any("filter", req.String()).Any("trace", logger.LogTraceData(r.Context())).Msg("List")
	out, err := s.app.List(r.Context(), req.ToFilter())
	if err != nil {
		web.SendError(w, err)
		return
	}
	web.SendPayload(w, "ok", out)
}

func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	var req peoplepb.GetRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(w, err)
	}
	if req.Id == "" {
		web.SendErrorBadRequest(w, "empty id", nil)
		return
	}
	logger.Log().Info().Any("id", &req.Id).Any("trace", logger.LogTraceData(r.Context())).Msg("Get")
	p, err := s.app.GetById(r.Context(), req.Id)
	if err != nil {
		logger.Log().Info().Any("id", &req.Id).Err(err).Any("trace", logger.LogTraceData(r.Context())).Msg("Get")
		web.SendError(w, err)
		return
	}
	web.SendPayload(w, "ok", p)
}

func (s *server) Save(w http.ResponseWriter, r *http.Request) {
	var req peoplepb.SaveRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(w, err)
	}
	p := req.ToPerson()
	logger.Log().Info().Any("person", &req).Any("trace", logger.LogTraceData(r.Context())).Msg("Save")
	err = s.app.Save(r.Context(), p)
	if err != nil {
		web.SendError(w, err)
		return
	}
	web.SendPayload(w, "ok", peoplepb.SaveResponse{Id: p.Id})
}

func (s *server) Update(w http.ResponseWriter, r *http.Request) {
	var req peoplepb.UpdateRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(w, err)
	}
	p := req.ToPerson()
	logger.Log().Info().Any("person", &req).Any("", logger.LogTraceData(r.Context())).Msg("Save")
	err = s.app.Save(r.Context(), p)
	if err != nil {
		web.SendError(w, err)
		return
	}
	web.SendPayload(w, "ok", nil)
}
