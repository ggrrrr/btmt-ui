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
		logger.ErrorCtx(r.Context(), err).Msg("List")
		web.SendError(w, err)
		return
	}
	logger.InfoCtx(r.Context()).Any("filter", req.String()).Msg("List")
	out, err := s.app.List(r.Context(), req.ToFilter())
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("List")
		web.SendError(w, err)
		return
	}
	people := []*peoplepb.Person{}
	for _, p := range out {
		person := peoplepb.FromPerson(&p)
		logger.DebugCtx(r.Context()).Any("person from", p).Msg("ListResult")
		logger.DebugCtx(r.Context()).Any("person to", person).Msg("ListResult")
		people = append(people, person)
	}
	web.SendPayload(w, "ok", people)
}

func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	var req peoplepb.GetRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("Get")
		web.SendError(w, err)
		return

	}
	if req.Id == "" {
		logger.ErrorCtx(r.Context(), err).Str("error", "empty id").Msg("Get")
		web.SendErrorBadRequest(w, "empty id", nil)
		return
	}
	logger.InfoCtx(r.Context()).Any("id", &req.Id).Msg("Get")
	p, err := s.app.GetById(r.Context(), req.Id)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("Get")
		web.SendError(w, err)
		return
	}
	web.SendPayload(w, "ok", p)
}

func (s *server) Save(w http.ResponseWriter, r *http.Request) {
	var req peoplepb.SaveRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("Save")
		web.SendError(w, err)
		return
	}
	p := req.ToPerson()
	logger.InfoCtx(r.Context()).Any("person", &req).Msg("Save")
	err = s.app.Save(r.Context(), p)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("Save")
		web.SendError(w, err)
		return
	}
	web.SendPayload(w, "ok", peoplepb.SaveResponse{Id: p.Id})
}

func (s *server) Update(w http.ResponseWriter, r *http.Request) {
	var req peoplepb.UpdateRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Any("person", &req).Msg("Update")
		web.SendError(w, err)
		return
	}
	p := req.ToPerson()
	logger.InfoCtx(r.Context()).Any("person", &req).Msg("Update")
	err = s.app.Save(r.Context(), p)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Any("person", &req).Msg("Update")
		web.SendError(w, err)
		return
	}
	web.SendPayload(w, "ok", nil)
}
