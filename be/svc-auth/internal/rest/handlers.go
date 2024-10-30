package rest

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
)

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

func (s *server) UserList(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.UserList", nil)
	defer func() {
		span.End(err)
	}()

	list, err := s.app.UserList(ctx)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("UserList")
		web.SendError(ctx, w, err)
		return
	}
	out := []*authpb.UserListPayload{}
	for _, a := range list {
		line := &authpb.UserListPayload{
			Username:    a.Subject,
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
	res, err := s.app.LoginPasswd(ctx, req.Username, req.Password)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("LoginPasswd")
		web.SendError(ctx, w, err)
		return
	}

	logger.InfoCtx(ctx).
		Str("Username", req.Username).
		Str("exp", res.AccessToken.ExpiresAt.String()).
		Msg("LoginPasswd")

	out := authpb.LoginTokenPayload{
		Username: req.Username,
		AccessToken: &authpb.LoginToken{
			Value:     res.AccessToken.Value,
			ExpiresAt: timestamppb.New(res.AccessToken.ExpiresAt),
		},
		RefreshToken: &authpb.LoginToken{
			Value:     res.RefreshToken.Value,
			ExpiresAt: timestamppb.New(res.RefreshToken.ExpiresAt),
		},
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

	err = s.app.TokenValidate(ctx)
	if err != nil {
		logger.ErrorCtx(r.Context(), err).Msg("TokenValidate")
		web.SendError(ctx, w, err)
		return
	}
	logger.InfoCtx(r.Context()).Msg("Validate")
	web.SendPayload(ctx, w, "ok", nil)
}

func (s *server) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := logger.Span(r.Context(), "rest.TokenRefresh", nil)
	defer func() {
		span.End(err)
	}()

	loginToken, err := s.app.TokenRefresh(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("TokenRefresh")
		web.SendError(ctx, w, err)
		return
	}

	out := authpb.LoginTokenPayload{
		Username: loginToken.Subject,
		AccessToken: &authpb.LoginToken{
			Value:     loginToken.AccessToken.Value,
			ExpiresAt: timestamppb.New(loginToken.AccessToken.ExpiresAt),
		},
	}

	logger.InfoCtx(r.Context()).Msg("TokenRefresh")
	web.SendPayload(ctx, w, "ok", &out)
}
