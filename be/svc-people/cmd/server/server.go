package server

import (
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	people "github.com/ggrrrr/btmt-ui/be/svc-people"
)

func Server() error {
	// var cfg auth.Cfg
	var cfg config.AppConfig

	err := config.InitConfig(&cfg)
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}

	err = people.Root(s.Waiter().Context(), s)
	if err != nil {
		return err
	}

	defer fmt.Println("module shutdown")

	s.Waiter().Add(
		s.WaitForWeb,
		s.WaitForGRPC,
	)

	return s.Waiter().Wait()
}
