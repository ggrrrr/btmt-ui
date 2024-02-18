package auth

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/app"
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

	logger.Info().Msg("starting...")
	return app.ErrTeapot
}
