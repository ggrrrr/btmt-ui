package grpc

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/authpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ authpb.AuthSvcServer = (*server)(nil)

func (*server) UserCreate(context.Context, *authpb.UserCreateRequest) (*authpb.UserCreateResponse, error) {

	panic("unimplemented")
}

func (s *server) UserList(ctx context.Context, _ *authpb.UserListRequest) (*authpb.UserListResponse, error) {
	logger.InfoCtx(ctx).Msg("UserList")
	list, err := s.app.UserList(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("UserList")
		return nil, app.ToGrpcError(err)
	}

	out := authpb.UserListResponse{
		Payload: []*authpb.UserListPayload{},
	}

	for _, a := range list.Payload() {
		out.Payload = append(out.Payload, &authpb.UserListPayload{
			Email:       a.Email,
			Status:      string(a.Status),
			SystemRoles: a.SystemRoles,
			CreatedAt:   timestamppb.New(a.CreatedAt),
		})
	}

	return &out, nil
}

func (*server) UserUpdate(context.Context, *authpb.UserUpdateRequest) (*authpb.UserUpdateResponse, error) {
	panic("unimplemented")
}

func (s *server) UserChangePasswd(ctx context.Context, req *authpb.UserChangePasswdRequest) (*authpb.UserChangePasswdResponse, error) {
	logger.InfoCtx(ctx).Msg("UserChangePasswd")
	err := s.app.UserChangePasswd(ctx, req.Email, req.Password, req.NewPassword)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("UserChangePasswd")
		return nil, app.ToGrpcError(err)
	}

	return &authpb.UserChangePasswdResponse{}, nil
}

func (s *server) LoginPasswd(ctx context.Context, req *authpb.LoginPasswdRequest) (*authpb.LoginPasswdResponse, error) {
	logger.InfoCtx(ctx).Str("email", req.Email).Msg("LoginPasswd")
	res, err := s.app.LoginPasswd(ctx, req.Email, req.Password)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("LoginPasswd")
		return nil, app.ToGrpcError(err)
	}
	return &authpb.LoginPasswdResponse{
		Payload: &authpb.LoginTokenPayload{
			Email: req.Email,
			Token: string(res.Payload()),
		},
	}, nil
}

func (s *server) TokenValidate(ctx context.Context, _ *authpb.TokenValidateRequest) (*authpb.TokenValidateResponse, error) {
	logger.InfoCtx(ctx).Msg("TokenValidate")
	err := s.app.TokenValidate(ctx)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("TokenValidate")
		return nil, app.ToGrpcError(err)
	}

	return &authpb.TokenValidateResponse{}, nil
}

func (s *server) LoginOauth2(ctx context.Context, _ *authpb.LoginOauth2Request) (*authpb.LoginOauth2Response, error) {
	logger.ErrorCtx(ctx, app.ErrTeapot).Msg("LoginOauth2")
	return nil, app.ToGrpcError(app.ErrTeapot)
}

func (s *server) Oauth2Config(ctx context.Context, _ *authpb.Oauth2ConfigRequest) (*authpb.Oauth2ConfigResponse, error) {
	out := authpb.Oauth2ConfigResponse{
		Payload: &authpb.Oauth2ConfigPayload{},
	}
	return &out, nil

}
