package events

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/svc-events/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-events/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, s *system.System) (err error) {
	return Root(ctx, s)
}

func Root(ctx context.Context, s *system.System) error {
	a, err := app.New()
	if err != nil {
		return err
	}
	restApp := rest.New(a)
	s.Mux().Mount("/rest", restApp.Router())

	if s.Mux() == nil {
		return fmt.Errorf("system.Mux is nil")
	}

	logger.Info().Msg("starting...")
	return nil
}
