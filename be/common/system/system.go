package system

import (
	"context"
	"errors"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi/v5"
	grpcRuntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/ggrrrr/btmt-ui/be/common/buildversion"
	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
)

type (
	System struct {
		cfg          config.AppConfig
		mux          *chi.Mux
		waiter       waiter.Waiter
		gateway      *grpcRuntime.ServeMux
		grpc         *grpc.Server
		aws          *session.Session
		verifier     token.Verifier
		buildVersion string
	}

	Service interface {
		// current config
		Config() config.AppConfig

		// http / rest router
		Mux() *chi.Mux

		// group waiter for shutdown/startup
		Waiter() waiter.Waiter

		// GRPC server
		RPC() *grpc.Server
	}
)

var _ (Service) = (*System)(nil)

func NewSystem(cfg config.AppConfig) (*System, error) {
	s := System{
		cfg:          cfg,
		buildVersion: buildversion.BuildVersion(),
	}
	logger.Info().
		Str("build.version", s.buildVersion).
		Int("max.procs", runtime.GOMAXPROCS(0)).
		Msg("system.init...")

	if cfg.Otel.Enabled {
		err := logger.ConfigureOtel(context.Background())
		if err != nil {
			return nil, err
		}
	}

	err := s.initJwt()
	if err != nil {
		return nil, err
	}
	s.initMux()
	s.initGRPC()
	s.initAws()
	s.waiter = waiter.New(waiter.CatchSignals())
	s.waiter.Cleanup(logger.Shutdown)

	return &s, nil
}

func (s *System) initJwt() error {
	if s.cfg.Jwt.UseMock == "mock" {
		s.verifier = token.NewVerifierMock()
		return nil
	}
	if s.cfg.Jwt.CrtFile == "" {
		return errors.New("CRT_FILE is empty")
	}
	ver, err := token.NewVerifier(s.cfg.Jwt.CrtFile)
	if err != nil {
		logger.Error(err).Str("crt_file", s.cfg.Jwt.CrtFile).Send()
		return err
	}
	s.verifier = ver
	return nil
}

func (s *System) initAws() {
	if s.cfg.Aws.Endpoint == "" {
		return
	}
	region := s.cfg.Aws.Region
	endpoint := s.cfg.Aws.Endpoint
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		logger.Error(err).Msg("initAws")
		return
	}
	logger.Info().Str("region", region).Str("endpoint", endpoint).Msg("initAws")
	s.aws = sess
}

func (s *System) Config() config.AppConfig {
	return s.cfg
}

func (s *System) RPC() *grpc.Server {
	return s.grpc
}

func (s *System) Mux() *chi.Mux {
	return s.mux
}

func (s *System) Gateway() *grpcRuntime.ServeMux {
	return s.gateway
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) Aws() *session.Session {
	return s.aws
}
