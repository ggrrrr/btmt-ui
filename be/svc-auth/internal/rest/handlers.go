package rest

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
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
	router.Get("/v1/noJson400", noJson400)
	router.Get("/v1/json500", json500)

	router.Post("/login/passwd", s.LoginPasswd)
	router.Get("/token/validate", s.TokenValidate)
	router.Post("/token/validate", s.TokenValidate)
	router.Post("/user/list", s.UserList)
	router.Get("/user/list", s.UserList)

	return router
}

func noJson400(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	// nolint: errcheck
	w.Write([]byte("noJson400"))
}

func json500(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	// nolint: errcheck
	w.Write([]byte(`{"me":"json500"}`))
}

func (s *server) LoginPasswd(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.LoginPasswd", nil)
	defer func() {
		span.End(err)
	}()

	var req authpb.LoginPasswdRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	logger.InfoCtx(r.Context()).Any("email", &req.Email).Msg("LoginPasswd")
	res, err := s.app.LoginPasswd(r.Context(), req.Email, req.Password)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("LoginPasswd")
		web.SendError(ctx, w, err)
		return
	}
	out := authpb.LoginTokenPayload{
		Email:     req.Email,
		Token:     res.Token,
		ExpiresAt: timestamppb.New(res.ExpiresAt),
	}

	// cookie := http.Cookie{
	// 	Name:     "accessToken",
	// 	Value:    string(res.Payload()),
	// 	HttpOnly: true,
	// 	Secure:   false,
	// 	Domain:   "localhost:8010",
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(365 * 24 * time.Hour),
	// }
	// http.SetCookie(w, &cookie)
	web.SendPayload(ctx, w, "ok", &out)
}

func (s *server) TokenValidate(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.TokenValidate", nil)
	defer func() {
		span.End(err)
	}()

	err = s.app.TokenValidate(r.Context())
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("TokenValidate")
		web.SendError(ctx, w, err)
		return
	}
	logger.InfoCtx(r.Context()).Msg("Validate")
	web.SendPayload(ctx, w, "ok", nil)
}

func (s *server) UserList(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.UserList", nil)
	defer func() {
		span.End(err)
	}()

	list, err := s.app.UserList(r.Context())
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("UserList")
		web.SendError(ctx, w, err)
		return
	}
	out := []*authpb.UserListPayload{}
	for _, a := range list {
		line := &authpb.UserListPayload{
			Email:       a.Email,
			Status:      string(a.Status),
			SystemRoles: a.SystemRoles,
			CreatedAt:   timestamppb.New(a.CreatedAt),
		}
		fmt.Printf("\t\t handler %v \n", a)
		if len(a.RealmRoles) > 0 {
			line.TenantRoles = map[string]*authpb.ListText{}
			for k := range a.RealmRoles {
				roles := authpb.ListText{
					List: a.RealmRoles[k],
				}
				line.TenantRoles[k] = &roles
			}
		}
		fmt.Printf("\t\t %v \n", line)
		out = append(out, line)
	}
	logger.InfoCtx(r.Context()).Msg("UserList")
	web.SendPayload(r.Context(), w, "ok", out)
}
