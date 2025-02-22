package main

import (
	"context"
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	"github.com/ggrrrr/btmt-ui/be/common/web"
	email "github.com/ggrrrr/btmt-ui/be/svc-email"
)

type appCfg struct {
	System system.Config
	WEB    web.Config
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	cfg := appCfg{}
	config.MustParse(&cfg)

	s, err := system.NewSystem(
		cfg.System,
		system.WithWebServer(cfg.WEB),
	)
	if err != nil {
		return err
	}

	m := email.Module{}

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
