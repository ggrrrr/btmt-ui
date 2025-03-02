package system

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"

	"github.com/ggrrrr/btmt-ui/be/common/buildversion"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
	"github.com/ggrrrr/btmt-ui/be/common/web"
)

type (
	SystemOptions func(s *System) error

	Config struct {
		AppName string `env:"APP_NAME"`
		OTEL    tracer.Config
		LOG     log.Config
		JWT     struct {
			UseMock string `env:"USE_MOCK"`
			Config  token.Config
		} `envPrefix:"JWT_"`
	}

	System struct {
		cfg       Config
		waiter    waiter.Waiter
		verifier  token.Verifier
		signer    token.Signer
		webServer *web.Server
	}

	Service interface {
		// group waiter for shutdown/startup
		Waiter() waiter.Waiter
		Verifier() token.Verifier
		Signer() token.Signer
		MountHandler(pattern string, router http.Handler) error
	}
)

var _ (Service) = (*System)(nil)

func NewSystem(cfg Config, opts ...SystemOptions) (*System, error) {
	s := &System{
		cfg: cfg,
	}
	log.Log().Info("build.version",
		slog.String("version", buildversion.BuildVersion()),
		slog.Int("max.procs", runtime.GOMAXPROCS(0)),
	)
	if cfg.OTEL.Client.Target != "" {
		err := tracer.Configure(context.Background(), cfg.AppName, cfg.OTEL)
		if err != nil {
			return nil, err
		}
	}

	err := s.initJWT()
	if err != nil {
		return nil, err
	}
	s.waiter = waiter.New(waiter.CatchSignals())
	s.waiter.AddCleanup(func() {
		tracer.Shutdown(context.Background())
	})

	for _, optFn := range opts {
		err = optFn(s)
		if err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *System) initJWT() error {
	if s.cfg.JWT.UseMock == "mock" {
		s.verifier = token.NewVerifierMock()
		s.signer = token.NewSignerMock()
		return nil
	}

	if s.cfg.JWT.Config.CrtFile == "" {
		return errors.New("CRT_FILE is empty")
	}

	ver, err := token.NewVerifier(s.cfg.JWT.Config.CrtFile)
	if err != nil {
		log.Log().Error(err, "NewVerifier")
		return err
	}
	s.verifier = ver

	if s.cfg.JWT.Config.KeyFile != "" {
		signer, err := token.NewSigner(s.cfg.JWT.Config.KeyFile)
		if err != nil {
			return err
		}
		s.signer = signer
	}

	return nil
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) Verifier() token.Verifier {
	return s.verifier
}

func (s *System) Signer() token.Signer {
	return s.signer
}

func (s *System) MountHandler(pattern string, router http.Handler) error {
	if s.webServer == nil {
		return fmt.Errorf("webServer is nil")
	}
	s.webServer.MountHandler(pattern, router)
	return nil
}

func WithWebServer(cfg web.Config) SystemOptions {
	return func(s *System) error {
		webServer, err := web.NewServer(
			"system",
			cfg,
			web.WithVerifier(s.verifier),
		)
		if err != nil {
			return err
		}
		s.webServer = webServer

		s.waiter.Add(func(ctx context.Context) error {
			return s.webServer.Startup()
		})

		s.waiter.AddCleanup(func() {
			s.webServer.Shutdown(context.Background())
		})

		return nil
	}
}
