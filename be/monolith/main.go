package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	auth "github.com/ggrrrr/btmt-ui/be/svc-auth"
	people "github.com/ggrrrr/btmt-ui/be/svc-people"
	tmpl "github.com/ggrrrr/btmt-ui/be/svc-tmpl"
)

type appCfg struct {
	System system.Config
	WEB    web.Config
}
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
	var cfg appCfg
	config.MustParse(&cfg)

	s, err := system.NewSystem(
		cfg.System,
		system.WithWebServer(cfg.WEB),
	)
	if err != nil {
		return err
	}

	m := &monolith{
		System: s,
		modules: []system.Module{
			&auth.Module{},
			&people.Module{},
			&tmpl.Module{},
		},
	}

	if err = m.configure(); err != nil {
		return err
	}
	if err = m.startup(); err != nil {
		return err
	}

	return m.Waiter().Wait()
}

func (m *monolith) configure() error {
	for i := range m.modules {
		ctx := m.Waiter().Context()
		if err := m.modules[i].Configure(ctx, m.System); err != nil {
			log.Log().Error(err, "fail", slog.String("module", m.modules[i].Name()))
			return err
		}
		log.Log().Info("configure", slog.String("module", m.modules[i].Name()))
	}

	return nil
}

func (m *monolith) startup() error {
	for i := range m.modules {
		ctx := m.Waiter().Context()
		if err := m.modules[i].Startup(ctx); err != nil {
			log.Log().Error(err, "fail", slog.String("module", m.modules[i].Name()))
			return err
		}
		log.Log().Info("startup", slog.String("module", m.modules[i].Name()))
	}

	return nil
}
