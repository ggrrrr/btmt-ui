package token

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

func (s *verifier) UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var peerInfo *peer.Peer
	var userRequest roles.UserRequest
	var authInfo roles.AuthInfo
	var ok bool
	var md metadata.MD
	var err error

	peerInfo, _ = peer.FromContext(ctx)

	startTs := time.Now()
	infoLog := logger.Info()
	defer func() {
		infoLog.TimeDiff("ts", time.Now(), startTs).Msg("UnaryInterceptor")
	}()

	if md, ok = metadata.FromIncomingContext(ctx); ok {
		userRequest = roles.FromGrpcMetadata(md, info.FullMethod)

		if userRequest.Authorization.AuthScheme != "" {
			authInfo, err = s.Verify(userRequest.Authorization)
			if err != nil {
				infoLog.
					Any("request", userRequest).
					Err(err)
				return req, status.Error(codes.Unauthenticated, err.Error())
			}
			infoLog.Str("user", authInfo.User)
		}
	}

	infoLog.
		Any("peerInfo", peerInfo).
		Any("device", userRequest.Device).
		Str("AuthScheme", userRequest.Authorization.AuthScheme).
		Str("FullMethod", userRequest.FullMethod)

	authInfo.Device = userRequest.Device
	ctx = roles.CtxWithAuthInfo(ctx, authInfo)

	return handler(ctx, req)
}
