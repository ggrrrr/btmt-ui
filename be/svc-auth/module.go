package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/postgres"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/app"
	repoPG "github.com/ggrrrr/btmt-ui/be/svc-auth/internal/repo/postgres"
	"github.com/ggrrrr/btmt-ui/be/svc-auth/internal/rest"
)

type (
	Config struct {
		DB       postgres.Config `envPrefix:"AUTH_PG_"`
		TokenTTL struct {
			AccessDuration  time.Duration `env:"ACCESS_DURATION" envDefault:"1m"`
			RefreshDuration time.Duration `env:"REFRESH_DURATION" envDefault:"30m"`
		} `envPrefix:"AUTH_"`
	}

	Module struct {
		cfg    Config
		app    app.App
		system system.Service
	}
)

var _ (system.Module) = (*Module)(nil)

func (m *Module) Module() app.App {
	return m.app
}

func (m *Module) Name() string {
	return "svc-auth"
}

func (m *Module) Configure(ctx context.Context, s system.Service) error {
	var err error
	cfg := Config{}
	config.MustParse(&cfg)

	m.cfg = cfg
	m.system = s

	db, err := postgres.Connect(cfg.DB)
	if err != nil {
		log.Log().Error(err, "cant connect to pg")
		return err
	}
	s.Waiter().AddCleanup(func() {
		db.Close()
	})

	pg, err := repoPG.Init(db)
	if err != nil {
		return err
	}

	m.app, err = app.New(
		app.WithAuthRepo(pg),
		app.WithHistoryRepo(pg),
		app.WithTokenSigner(s.Signer()),
		app.WithTokenTTL(m.cfg.TokenTTL.AccessDuration, m.cfg.TokenTTL.RefreshDuration),
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) Startup(ctx context.Context) (err error) {
	if m.app == nil {
		log.Log().Error(fmt.Errorf("app is nil"), "starutup")
		panic("m.app is nil")
	}
	return m.system.MountHandler("/v1/auth", rest.Router(rest.Handler(m.app)))
}

// func initAwsRepo(_ context.Context, _ waiter.Waiter, cfg config.AppConfig) (ddd.AuthRepo, error) {
// 	repo, err := dynamodb.New(cfg.Aws, cfg.Dynamodb)
// 	if err != nil {
// 		logger.Error(err).Msg("initAwsRepo error")
// 		return nil, err
// 	}
// 	// s.Waiter().Cleanup(func() {
// 	// 	repoDb.Close(ctx)
// 	// })
// 	return repo, nil
// }
