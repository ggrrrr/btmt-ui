package system

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

func (s *Server) unaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var peerInfo *peer.Peer
	var userRequest app.RequestIn
	var authInfo roles.AuthInfo
	var ok bool
	var md metadata.MD
	var err error

	peerInfo, _ = peer.FromContext(ctx)

	if s.verifier == nil {
		log.Log().WarnCtx(ctx, err, "verifier is nil",
			slog.String("rpc.server", s.name))

		return handler(ctx, req)
	}

	startTs := time.Now()
	slogAttr := []slog.Attr{
		slog.String("rpc.server", s.name),
	}
	defer func() {
		slogAttr = append(slogAttr,
			slog.Duration("ts", time.Now().Sub(startTs)),
		)
		log.Log().Info("unaryInterceptor", slogAttr...)
	}()

	if md, ok = metadata.FromIncomingContext(ctx); ok {
		userRequest = roles.FromGrpcMetadata(md, info.FullMethod)

		if !userRequest.AuthData.IsZero() {
			authInfo, err = s.verifier.Verify(userRequest.AuthData)
			if err != nil {
				slogAttr = append(slogAttr,
					slog.Any("request", userRequest),
					slog.Any("error", err),
				)
				return req, status.Error(codes.Unauthenticated, err.Error())
			}
			slogAttr = append(slogAttr,
				slog.String("subject", authInfo.Subject),
			)
		}
	}

	slogAttr = append(slogAttr,
		slog.Any("peerInfo", peerInfo),
		slog.Any("device", userRequest.Device),
		slog.String("AuthScheme", userRequest.AuthData.AuthScheme),
		slog.String("FullMethod", userRequest.FullMethod),
	)

	authInfo.Device = userRequest.Device
	ctx = roles.CtxWithAuthInfo(ctx, authInfo)

	return handler(ctx, req)
}
