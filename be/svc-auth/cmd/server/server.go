package server

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	auth "github.com/ggrrrr/btmt-ui/be/svc-auth"
)

type appCfg struct {
	System system.Config
	WEB    web.Config
	OTEL   logger.Config
}

func Server() error {
	cfg := appCfg{}
	config.MustParse(&cfg)

	s, err := system.NewSystem(
		cfg.System,
		system.WithWebServer(cfg.WEB),
	)
	if err != nil {
		return err
	}

	m := auth.Module{}

	err = m.Configure(context.Background(), s)
	if err != nil {
		return err
	}

	err = m.Startup(context.Background())
	if err != nil {
		return err
	}

	defer fmt.Println("module shutdown")

	return s.Waiter().Wait()
}
