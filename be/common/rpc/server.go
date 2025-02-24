package system

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	ServerOption func(*Server) error

	Config struct {
		ListenAddress string `env:"RPC_LISTEN_ADDRESS"`
	}

	RPCServer interface {
		Startup(ctx context.Context) error
		Shutdown(ctx context.Context) error
	}

	Server struct {
		grpcServer *grpc.Server
		name       string
		cfg        Config
		verifier   token.Verifier
	}
)

var _ (RPCServer) = (*Server)(nil)

func NewServer(name string, cfg Config, opts ...ServerOption) (*Server, error) {
	if cfg.ListenAddress == "" {
		return nil, fmt.Errorf("empty LISTEN_ADDRESS")
	}
	out := &Server{
		name: name,
		cfg:  cfg,
	}

	out.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(out.unaryInterceptor))

	return out, nil
}

func (s *Server) Startup(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.ListenAddress)
	if err != nil {
		logger.Error(err).
			Str("rpc.server", s.name).
			Str("address", s.cfg.ListenAddress).
			Msg("Failed to listen")
		return err
	}

	err = s.grpcServer.Serve(listener)
	if err != nil && err != grpc.ErrServerStopped {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

func WithVerifier(v token.Verifier) ServerOption {
	return func(s *Server) error {
		s.verifier = v
		return nil
	}
}
