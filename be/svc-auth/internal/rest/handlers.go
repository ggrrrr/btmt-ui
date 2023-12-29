package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
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
	router.Post("/v1/noJson400", noJson400)
	router.Post("/v1/json500", json500)

	router.Post("/v1/auth/login/passwd", s.LoginPasswd)
	router.Get("/v1/auth/validate", s.TokenValidate)
	router.Post("/v1/auth/validate", s.TokenValidate)
	router.Post("/v1/auth/user/list", s.UserList)

	return router
}

func noJson400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	w.Write([]byte("asdasd"))
}

func json500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte(`{"me":"json500"}`))
}

func (s *server) LoginPasswd(w http.ResponseWriter, r *http.Request) {
	var req authpb.LoginPasswdRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(w, err)
		return
	}
	logger.InfoCtx(r.Context()).Any("email", &req.Email).Msg("LoginPasswd")
	res, err := s.app.LoginPasswd(r.Context(), req.Email, req.Password)
	if err != nil {
		fmt.Printf("%+v \n", err)
		logger.ErrorCtx(r.Context(), err).Msg("LoginPasswd")
		web.SendError(w, err)
		return
	}
	out := authpb.LoginTokenPayload{
		Email: req.Email,
		Token: string(res.Payload()),
	}
	web.SendPayload(w, "ok", &out)
}

func (s *server) TokenValidate(w http.ResponseWriter, r *http.Request) {
	err := s.app.TokenValidate(r.Context())
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("TokenValidate")
		web.SendError(w, err)
		return
	}
	logger.InfoCtx(r.Context()).Msg("Validate")
	web.SendPayload(w, "ok", nil)
}

func (s *server) UserList(w http.ResponseWriter, r *http.Request) {
	list, err := s.app.UserList(r.Context())
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("UserList")
		web.SendError(w, err)
		return
	}
	out := []authpb.UserListPayload{}
	for _, a := range list.Payload() {
		out = append(out, authpb.UserListPayload{
			Email:       a.Email,
			Status:      string(a.Status),
			SystemRoles: a.SystemRoles,
			CreatedAt:   timestamppb.New(a.CreatedAt),
		})
	}
	logger.InfoCtx(r.Context()).Msg("UserList")
	web.SendPayload(w, "ok", out)
}
