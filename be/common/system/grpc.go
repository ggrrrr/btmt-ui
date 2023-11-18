package system

import (
	"context"
	"fmt"
	"net"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

func (s *System) WaitForGRPC(ctx context.Context) error {
	if s.cfg.Grpc.Address == "" {
		logger.Info().Msg("Address is empty, initGRPC skip")
		return nil
	}

	addr := s.cfg.Grpc.Address
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error(err).Str("address", addr).Msg("Failed to listen")
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		logger.Info().Str("address", addr).Msg("grpc started")
		defer fmt.Println("rpc server shutdown")
		if err := s.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		logger.Info().Str("address", addr).Msg("grpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(time.Duration(s.cfg.ShutdownTimeout))
		select {
		case <-timeout.C:
			// Force it to stop
			s.RPC().Stop()
			return fmt.Errorf("grpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func (s *System) initGRPC() {
	if s.cfg.Grpc.Address == "" {
		logger.Debug().Msg("initGRPC skip")
		return
	}
	s.grpc = grpc.NewServer(
		grpc.UnaryInterceptor(s.unaryInterceptor),
	)
}

func (s *System) unaryInterceptor(
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
		infoLog.TimeDiff("ts", time.Now(), startTs).Msg("unaryInterceptor")
	}()

	if md, ok = metadata.FromIncomingContext(ctx); ok {
		userRequest = roles.FromGrpcMetadata(md, info.FullMethod)

		if userRequest.Authorization.AuthScheme != "" {
			authInfo, err = s.verifier.Verify(userRequest.Authorization)
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
