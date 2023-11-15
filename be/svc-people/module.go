package people

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/mongodb"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/grpc"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, s *system.System) (err error) {
	return Root(ctx, s)
}

func Root(ctx context.Context, s *system.System) error {

	repoDb, err := mongodb.New(ctx, s.Config().Mongo)
	if err != nil {
		logger.Log().Error().Err(err).Msg("db")
		return err
	}

	s.Waiter().Cleanup(func() {
		repoDb.Close(ctx)
	})

	appRepo := repo.New(s.Config().Mongo.Collection, repoDb)

	a, err := app.New(
		app.WithPeopleRepo(appRepo),
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

	grpc.RegisterServer(a, s.RPC())

	logger.Log().Info().Msg("starting...")
	if err = rest.RegisterGateway(ctx, s.Gateway(), "localhost:8081"); err != nil {
		return err
	}
	return nil
}
