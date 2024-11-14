package system

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/common/ver"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
)

type (
	System struct {
		cfg          config.AppConfig
		mux          *chi.Mux
		gateway      *runtime.ServeMux
		waiter       waiter.Waiter
		grpc         *grpc.Server
		aws          *session.Session
		verifier     token.Verifier
		buildVersion string
		buildTime    time.Time
	}

	Service interface {
		Config() config.AppConfig
		Mux() *chi.Mux
		Waiter() waiter.Waiter
		RPC() *grpc.Server
	}
)

var _ (Service) = (*System)(nil)

func NewSystem(cfg config.AppConfig) (*System, error) {
	s := System{
		cfg:          cfg,
		buildVersion: ver.BuildVersion(),
		buildTime:    ver.BuildTime(),
	}
	logger.Info().
		Str("build.version", s.buildVersion).
		Time("build.time", s.buildTime).
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

func (s *System) Gateway() *runtime.ServeMux {
	return s.gateway
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) Aws() *session.Session {
	return s.aws
}
