package main

import (
	"fmt"
	"os"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	auth "github.com/ggrrrr/btmt-ui/be/svc-auth"
	people "github.com/ggrrrr/btmt-ui/be/svc-people"
	tmpl "github.com/ggrrrr/btmt-ui/be/svc-tmpl"
)

type monolith struct {
	*system.System
	modules []system.Module
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	var err error
	var cfg config.AppConfig
	err = config.InitConfig(&cfg)
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	m := monolith{
		System: s,
		modules: []system.Module{
			&auth.Module{},
			&people.Module{},
			&tmpl.Module{},
		},
	}

	if err = m.startupModules(); err != nil {
		return err
	}

	m.Waiter().Add(
		m.WaitForWeb,
		m.WaitForGRPC,
	)

	return m.Waiter().Wait()
}

func (m *monolith) startupModules() error {
	for _, module := range m.modules {
		ctx := m.Waiter().Context()
		if err := module.Startup(ctx, m); err != nil {
			logger.Error(err).Str("module", module.Name()).Msg("failed")
			return err
		}
		logger.Info().Str("module", module.Name()).Msg("added")
	}

	return nil
}
