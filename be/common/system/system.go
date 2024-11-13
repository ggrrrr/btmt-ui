package system

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/common/ver"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
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
