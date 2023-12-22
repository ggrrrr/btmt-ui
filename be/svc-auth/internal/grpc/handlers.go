package grpc

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
)

var _ authpb.AuthSvcServer = (*server)(nil)

func (*server) CreateAuth(context.Context, *authpb.CreateAuthRequest) (*authpb.CreateAuthResponse, error) {

	panic("unimplemented")
}

func (*server) UpdateAuth(context.Context, *authpb.UpdateAuthRequest) (*authpb.UpdateAuthResponse, error) {
	panic("unimplemented")
}

func (s *server) LoginPasswd(ctx context.Context, req *authpb.LoginPasswdRequest) (*authpb.LoginPasswdResponse, error) {
	logger.InfoCtx(ctx).Str("email", req.Email).Msg("LoginPasswd")
	res, err := s.app.LoginPasswd(ctx, req.Email, req.Password)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("LoginPasswd")
		return nil, app.ToGrpcError(err)
	}
	return &authpb.LoginPasswdResponse{
		Payload: &authpb.LoginPasswdPayload{
			Email: req.Email,
			Token: string(res.Payload()),
		},
	}, nil
}

func (s *server) ChangePasswd(ctx context.Context, req *authpb.ChangePasswdRequest) (*authpb.ChangePasswdResponse, error) {
	logger.InfoCtx(ctx).Msg("ChangePasswd")
	err := s.app.ChangePasswd(ctx, req.Email, req.Password, req.NewPassword)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ChangePasswd")
		return nil, app.ToGrpcError(err)
	}

	return &authpb.ChangePasswdResponse{}, nil
}

func (s *server) ValidateToken(ctx context.Context, _ *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	logger.InfoCtx(ctx).Msg("ValidateToken")
	err := s.app.Validate(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ValidateToken")
		return nil, app.ToGrpcError(err)
	}

	return &authpb.ValidateTokenResponse{}, nil
}

func (s *server) ListAuth(ctx context.Context, _ *authpb.ListAuthRequest) (*authpb.ListAuthResponse, error) {
	logger.InfoCtx(ctx).Msg("ListAuth")
	list, err := s.app.ListAuth(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("ListAuth")
		return nil, app.ToGrpcError(err)
	}

	out := authpb.ListAuthResponse{
		Payload: []*authpb.ListAuthPayload{},
	}

	for _, a := range list.Payload() {
		out.Payload = append(out.Payload, &authpb.ListAuthPayload{
			Email:       a.Email,
			Status:      string(a.Status),
			SystemRoles: a.SystemRoles,
			CreatedAt:   a.CreatedAt.GoString(),
		})
	}

	return &out, nil
}

func (s *server) LoginOauth2(ctx context.Context, _ *authpb.LoginOauth2Request) (*authpb.LoginOauth2Response, error) {
	logger.ErrorCtx(ctx, fmt.Errorf("ErrTeepot")).Msg("LoginOauth2")
	return nil, app.ToGrpcError(app.ErrTeepot)
}

func (s *server) Oauth2Config(ctx context.Context, _ *authpb.Oauth2ConfigRequest) (*authpb.Oauth2ConfigResponse, error) {
	out := authpb.Oauth2ConfigResponse{
		Payload: &authpb.Oauth2ConfigPayload{},
	}
	return &out, nil

}
