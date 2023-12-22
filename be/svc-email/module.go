package auth

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	// "github.com/ggrrrr/btmt-ui/be/svc-auth/internal/grpc"
	// "github.com/ggrrrr/btmt-ui/be/svc-auth/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, s *system.System) (err error) {
	return Root(ctx, s)
}

func Root(ctx context.Context, s *system.System) error {
	// a
	restApp := rest.New()
	s.Mux().Mount("/rest", restApp.Router())

	if s.Mux() == nil {
		return fmt.Errorf("system.Mux is nil")
	}

	if s.Config().Grpc.Address != "" {
		grpc.RegisterServer(a, s.RPC())
		if err = rest.RegisterGateway(ctx, s.Gateway(), "localhost:8011"); err != nil {
			return err
		}
	}

	logger.Info().Msg("starting...")
	return nil
}
