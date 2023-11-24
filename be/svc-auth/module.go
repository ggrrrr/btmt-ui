package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/ddd"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/grpc"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/dynamodb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/postgres"
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

func InitApp(ctx context.Context, w waiter.Waiter, cfg config.AppConfig) (app.App, error) {
	var err error
	var repo ddd.AuthPasswdRepo
	awsCfg := false
	pgCfg := false

	if cfg.Aws.Region != "" {
		awsCfg = true
	}

	if cfg.Postgres.Host != "" {
		pgCfg = true
	}

	if awsCfg == pgCfg {
		logger.Error(err).Any("pg", cfg.Postgres).Msg("repo init")
		return nil, errors.New(" postgres and aws dynamodb")
	}
	logger.Error(err).Any("aws", awsCfg).Any("pg", pgCfg).Msg("repo init")
	if awsCfg {
		repo, err = initAwsRepo(ctx, w, cfg)
		if err != nil {
			return nil, err
		}
	}
	if pgCfg {
		repo, err = initPgRepo(ctx, w, cfg)
		if err != nil {
			return nil, err
		}
	}
	if repo == nil {
		logger.Error(err).Msg("repo init")
		return nil, errors.New("shit")
	}

	tokemSigner, err := token.NewSigner(cfg.Jwt.TTL, cfg.Jwt.KeyFile)
	if err != nil {
		logger.Error(err).Msg("NewSigner")
		return nil, err
	}

	a, err := app.New(
		app.WithAuthRepo(repo),
		app.WithTokenSigner(tokemSigner),
	)
	if err != nil {
		logger.Error(err).Msg("app error")
		return nil, err
	}
	return a, nil
}

func initAwsRepo(ctx context.Context, w waiter.Waiter, cfg config.AppConfig) (ddd.AuthPasswdRepo, error) {
	repo, err := dynamodb.New(cfg.Aws)
	if err != nil {
		logger.Error(err).Msg("initAwsRepo error")
		return nil, err
	}
	// s.Waiter().Cleanup(func() {
	// 	repoDb.Close(ctx)
	// })
	return repo, nil
}

func initPgRepo(ctx context.Context, w waiter.Waiter, cfg config.AppConfig) (ddd.AuthPasswdRepo, error) {
	repo, err := postgres.Connect(cfg.Postgres)
	if err != nil {
		logger.Error(err).Msg("initPgRepo error")
		return nil, err
	}
	w.Cleanup(func() {
		repo.Close()
	})
	return repo, nil
}

func Root(ctx context.Context, s *system.System) error {
	a, err := InitApp(ctx, s.Waiter(), s.Config())
	if err != nil {
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

	logger.Info().Msg("starting...")
	return nil
}
