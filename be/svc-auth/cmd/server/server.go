package server

import (
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	auth "github.com/ggrrrr/btmt-ui/be/svc-auth"
)

func Server() error {
	// var cfg auth.Cfg
	var cfg config.AppConfig

	err := config.InitConfig(&cfg)
	logger.Init(cfg.Log)
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	err = auth.Root(s.Waiter().Context(), s)
	if err != nil {
		return err
	}

	defer fmt.Println("auth module shutdown")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForGRPC,
	)

	return s.Waiter().Wait()
}
