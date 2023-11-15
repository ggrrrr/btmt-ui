package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

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
	router.Post("/v1/nojson", noJson400)
	router.Post("/v1/auth/login/passwd", s.LoginPasswd)
	router.Get("/v1/auth/validate", s.Validate)
	router.Post("/v1/auth/validate", s.Validate)
	return router
}

func noJson400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	w.Write([]byte("asdasd"))
}

func (s *server) LoginPasswd(w http.ResponseWriter, r *http.Request) {
	var req authpb.LoginPasswdRequest
	err := web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(w, err)
		return
	}
	logger.Log().Info().Any("asd", &req.Password).Send()
	res, err := s.app.LoginPasswd(r.Context(), req.Email, req.Password)
	if err != nil {
		fmt.Printf("%+v \n", err)
		logger.Log().Error().Err(err).Send()
		web.SendError(w, err)
		return
	}
	out := authpb.LoginPasswdResponse{
		Email: req.Email,
		Token: string(res.Payload()),
	}
	web.SendPayload(w, "ok", &out)
}

func (s *server) Validate(w http.ResponseWriter, r *http.Request) {
	err := s.app.Validate(r.Context())
	if err != nil {
		logger.Log().Err(err).Send()
		web.SendError(w, err)
		return
	}
	logger.Log().Info().Msg("Validate")
	web.SendPayload(w, "ok", nil)
}
