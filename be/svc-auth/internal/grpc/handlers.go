package grpc

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
)

var _ authpb.AuthSvcServer = (*server)(nil)

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

func (s *server) UpdatePasswd(ctx context.Context, req *authpb.UpdatePasswdRequest) (*authpb.UpdatePasswdResponse, error) {
	logger.InfoCtx(ctx).Msg("UpdatePasswd")
	err := s.app.UpdatePasswd(ctx, req.Email, req.Password, req.NewPassword)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("UpdatePasswd")
		return nil, app.ToGrpcError(err)
	}

	return &authpb.UpdatePasswdResponse{}, nil
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

func (s *server) LoginOauth2(ctx context.Context, _ *authpb.LoginOauth2Request) (*authpb.LoginOauth2Response, error) {
	logger.ErrorCtx(ctx, fmt.Errorf("ErrTeepot")).Msg("LoginOauth2")
	return nil, app.ToGrpcError(app.ErrTeepot)
}

func (s *server) GetOauth2Config(ctx context.Context, _ *authpb.GetOauth2ConfigRequest) (*authpb.GetOauth2ConfigResponse, error) {
	logger.ErrorCtx(ctx, fmt.Errorf("ErrTeepot")).Msg("GetOauth2Config")
	out := authpb.GetOauth2ConfigResponse{
		Payload: &authpb.Oauth2ConfigPayload{},
	}
	return &out, app.ToGrpcError(app.ErrTeepot)

}
