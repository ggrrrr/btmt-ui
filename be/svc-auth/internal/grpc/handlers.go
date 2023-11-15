package grpc

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
)

var _ authpb.AuthSvcServer = (*server)(nil)

func (s *server) LoginPasswd(ctx context.Context, req *authpb.LoginPasswdRequest) (*authpb.LoginPasswdResponse, error) {
	logger.Log().Info().Str("email", req.Email).Any("info", logger.LogTraceData(ctx)).Msg("LoginPasswd")
	res, err := s.app.LoginPasswd(ctx, req.Email, req.Password)
	if err != nil {
		logger.Log().Error().Errs("LoginPasswd", []error{err}).Send()
		return nil, app.ToGrpcError(err)
	}
	return &authpb.LoginPasswdResponse{
		Email: req.Email,
		Token: string(res.Payload()),
	}, nil
}

func (s *server) UpdatePasswd(ctx context.Context, req *authpb.UpdatePasswdRequest) (*authpb.UpdatePasswdResponse, error) {
	logger.Log().Info().Any("info", logger.LogTraceData(ctx)).Msg("UpdatePasswd")
	err := s.app.UpdatePasswd(ctx, req.Email, req.Password, req.NewPassword)
	if err != nil {
		logger.Log().Error().Errs("UpdatePasswd", []error{err}).Send()
		return nil, app.ToGrpcError(err)
	}

	return &authpb.UpdatePasswdResponse{}, nil
}

func (s *server) ValidateToken(ctx context.Context, _ *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	logger.Log().Info().Any("info", logger.LogTraceData(ctx)).Msg("ValidateToken")
	err := s.app.Validate(ctx)
	if err != nil {
		logger.Log().Error().Errs("Validate", []error{err}).Send()
		return nil, app.ToGrpcError(err)
	}

	return &authpb.ValidateTokenResponse{}, nil
}

func (s *server) LoginOauth2(context.Context, *authpb.LoginOauth2Request) (*authpb.LoginOauth2Response, error) {
	return nil, app.ToGrpcError(app.ErrTeepot)
}
