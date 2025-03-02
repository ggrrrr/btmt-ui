package rest

import (
	"fmt"
	"net/http"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/web"
	authpb "github.com/ggrrrr/btmt-ui/be/svc-auth/authpb/v1"
)

var _ (AppHandler) = (*server)(nil)

func (s *server) UserList(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.UserList")
	defer func() {
		span.End(err)
	}()

	list, err := s.app.UserList(ctx)
	if err != nil {
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
	web.SendJSONPayload(r.Context(), w, "ok", out)
}

func (s *server) LoginPasswd(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.LoginPasswd")
	defer func() {
		span.End(err)
	}()

	var req authpb.LoginPasswdRequest
	err = web.DecodeJsonRequest(r, &req)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	loginPasswd, err := s.app.LoginPasswd(ctx, req.Username, req.Password)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	out := authpb.LoginTokenPayload{
		Username:      req.Username,
		AdminUsername: loginPasswd.AdminSubject,
		AccessToken: &authpb.LoginToken{
			Value:     loginPasswd.AccessToken.Value,
			ExpiresAt: timestamppb.New(loginPasswd.AccessToken.ExpiresAt),
		},
		RefreshToken: &authpb.LoginToken{
			Value:     loginPasswd.RefreshToken.Value,
			ExpiresAt: timestamppb.New(loginPasswd.RefreshToken.ExpiresAt),
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
	web.SendJSONPayload(ctx, w, "ok", &out)
}

func (s *server) TokenValidate(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.TokenValidate")
	defer func() {
		span.End(err)
	}()

	err = s.app.TokenValidate(ctx)
	if err != nil {
		web.SendError(ctx, w, err)
		return
	}
	web.SendJSONPayload(ctx, w, "ok", nil)
}

func (s *server) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx, span := s.tracer.Span(r.Context(), "rest.TokenValidate")

	defer func() {
		span.End(err)
	}()

	loginToken, err := s.app.TokenRefresh(ctx)
	if err != nil {
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

	web.SendJSONPayload(ctx, w, "ok", &out)
}
