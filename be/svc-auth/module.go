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
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/dynamodb"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/mem"
	repoPg "github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/postgres"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/rest"
)

type Module struct{}

var _ (system.Module) = (*Module)(nil)

func (m *Module) Name() string {
	return "svc-auth"
}

func (*Module) Startup(ctx context.Context, s system.Service) (err error) {
	return Root(ctx, s)
}

// Config arr sub services needed for the auth
func InitApp(ctx context.Context, s system.Service) (*app.Application, error) {
	var err error
	var repo ddd.AuthPasswdRepo
	awsCfg := false
	pgCfg := false
	logger.Info().
		Str("Aws.Region", s.Config().Aws.Region).
		Str("Postgres.Host", s.Config().Postgres.Host).
		Str("Postgres.User", s.Config().Postgres.Username).
		Send()

	if s.Config().Dynamodb.Database != "" {
		awsCfg = true
	}

	if s.Config().Postgres.Host != "" {
		pgCfg = true
	}

	if awsCfg == pgCfg {
		logger.Warn().Msg("in memory repo")
		pass, _ := app.HashPassword("asdasd")
		repo, _ = mem.New()
		asdUser := ddd.AuthPasswd{
			Email:       "asd@asd",
			Status:      ddd.StatusEnabled,
			SystemRoles: []string{"admin"},
			Passwd:      pass,
		}
		err = repo.Save(ctx, asdUser)
		if err != nil {
			logger.Error(err).Msg("InitApp.save")
		}
	}
	if awsCfg {
		repo, err = initAwsRepo(ctx, s.Waiter(), s.Config())
		if err != nil {
			return nil, err
		}
	}
	if pgCfg {
		repo, err = initPgRepo(s)
		if err != nil {
			return nil, err
		}
	}
	if repo == nil {
		logger.Error(err).Msg("repo init")
		return nil, errors.New("no repo")
	}

	tokenSigner, err := token.NewSigner(s.Config().Jwt.TTL, s.Config().Jwt.KeyFile)
	if err != nil {
		logger.Error(err).Msg("NewSigner")
		return nil, err
	}

	a, err := app.New(
		app.WithAuthRepo(repo),
		app.WithTokenSigner(tokenSigner),
	)
	if err != nil {
		logger.Error(err).Msg("app error")
		return nil, err
	}

	return a, nil
}

func initAwsRepo(ctx context.Context, w waiter.Waiter, cfg config.AppConfig) (ddd.AuthPasswdRepo, error) {
	repo, err := dynamodb.New(cfg.Aws, cfg.Dynamodb)
	if err != nil {
		logger.Error(err).Msg("initAwsRepo error")
		return nil, err
	}
	// s.Waiter().Cleanup(func() {
	// 	repoDb.Close(ctx)
	// })
	return repo, nil
}

func initPgRepo(s system.Service) (ddd.AuthPasswdRepo, error) {
	repo, err := repoPg.Init(s.DB())
	if err != nil {
		logger.Error(err).Msg("initPgRepo error")
		return nil, err
	}

	return repo, nil
}

func Root(ctx context.Context, s system.Service) error {
	logger.Info().Msg("Root")
	a, err := InitApp(ctx, s)
	if err != nil {
		return err
	}

	if s.Mux() == nil {
		return fmt.Errorf("system.Mux is nil")
	}
	restApp := rest.New(a)
	s.Mux().Mount("/auth", restApp.Router())

	// if s.Config().Grpc.Address != "" {
	// 	grpc.RegisterServer(a, s.RPC())
	// 	if err = rest.RegisterGateway(ctx, s.Gateway(), "localhost:8011"); err != nil {
	// 		return err
	// 	}
	// }

	// logger.Info().Msg("started")
	return nil
}
