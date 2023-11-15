package auth

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/grpc"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/dynamodb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/rest"
)

type Module struct{}

type Cfg struct {
	Prefix           string
	config.AppConfig `mapstructure:",squash"`
}

func (Module) Startup(ctx context.Context, s *system.System) (err error) {
	return Root(ctx, s)
}

func Root(ctx context.Context, s *system.System) error {
	repo, err := dynamodb.New(s.Aws(), s.Config().Aws.Database)

	if err != nil {
		logger.Log().Error().Err(err).Msg("awsdyno aws error")
		return err
	}

	tokemSigner, err := token.NewSigner(s.Config().Jwt.TTL, s.Config().Jwt.KeyFile)
	if err != nil {
		logger.Log().Error().Err(err).Msg("NewSigner")
		return err
	}

	a, err := app.New(
		app.WithAuthRepo(repo),
		app.WithTokenSigner(tokemSigner),
	)
	if err != nil {
		logger.Log().Error().Err(err).Msg("app error")
		return err
	}
	restApp := rest.New(a)
	s.Mux().Mount("/rest", restApp.Router())

	if s.Mux() == nil {
		return fmt.Errorf("system.Mux is nil")
	}

	if s.Config().Grpc.Address != "" {
		grpc.RegisterServer(a, s.RPC())
		if err = rest.RegisterGateway(ctx, s.Gateway(), "localhost:8081"); err != nil {
			return err
		}
	}

	logger.Log().Info().Msg("starting...")
	return nil
}
