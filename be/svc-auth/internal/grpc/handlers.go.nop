package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	authpb "github.com/ggrrrr/btmt-ui/be/svc-auth/authpb/v1"
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

	for _, a := range list {
		out.Payload = append(out.Payload, &authpb.UserListPayload{
			Username:    a.Subject,
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
	err := s.app.UserChangePasswd(ctx, req.Username, req.Password, req.NewPassword)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("UserChangePasswd")
		return nil, app.ToGrpcError(err)
	}

	return &authpb.UserChangePasswdResponse{}, nil
}

func (s *server) LoginPasswd(ctx context.Context, req *authpb.LoginPasswdRequest) (*authpb.LoginPasswdResponse, error) {
	logger.InfoCtx(ctx).Str("Username", req.Username).Msg("LoginPasswd")
	res, err := s.app.LoginPasswd(ctx, req.Username, req.Password)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("LoginPasswd")
		return nil, app.ToGrpcError(err)
	}
	return &authpb.LoginPasswdResponse{
		Payload: &authpb.LoginTokenPayload{
			Username: req.Username,
			AccessToken: &authpb.LoginToken{
				Value:     string(res.AccessToken.Value),
				ExpiresAt: timestamppb.New(res.AccessToken.ExpiresAt),
			},
			RefreshToken: &authpb.LoginToken{
				Value:     string(res.RefreshToken.Value),
				ExpiresAt: timestamppb.New(res.RefreshToken.ExpiresAt),
			},
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
