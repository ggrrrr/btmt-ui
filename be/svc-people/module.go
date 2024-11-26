package people

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/mgo"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/grpc"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/repo"
	"github.com/ggrrrr/btmt-ui/be/svc-people/internal/rest"
)

type Module struct{}

var _ (system.Module) = (*Module)(nil)

func (*Module) Name() string {
	return "svc-people"
}

func (*Module) Startup(ctx context.Context, s system.Service) (err error) {
	return Root(ctx, s)
}

func InitApp(ctx context.Context, cfg config.AppConfig) (*app.App, []waiter.CleanupFunc, error) {
	closeFns := []waiter.CleanupFunc{}
	db, err := mgo.New(ctx, cfg.Mgo)
	if err != nil {
		logger.Error(err).Msg("db")
		return nil, closeFns, err
	}
	fn := func() {
		db.Close(ctx)
	}
	closeFns = append(closeFns, fn)

	appRepo := repo.New(cfg.Mgo.Collection, db)
	a, err := app.New(
		app.WithPeopleRepo(appRepo),
	)
	if err != nil {
		logger.Error(err).Msg("app error")
		return nil, closeFns, err
	}
	return a, closeFns, nil
}

func Root(ctx context.Context, s system.Service) error {
	logger.Info().Msg("Root")
	a, fns, err := InitApp(ctx, s.Config())
	s.Waiter().Cleanup(fns...)
	if err != nil {
		logger.Error(err).Msg("app error")
		return err
	}

	restApp := rest.New(a)
	s.Mux().Mount("/people", restApp.Router())

	if s.Mux() == nil {
		return fmt.Errorf("system.Mux is nil")
	}

	grpc.RegisterServer(a, s.RPC())

	// logger.Info().Msg("starting...")
	// if err = rest.RegisterGateway(ctx, s.Gateway(), "localhost:8021"); err != nil {
	// 	return err
	// }
	return nil
}
