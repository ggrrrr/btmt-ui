package main

import (
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	email "github.com/ggrrrr/btmt-ui/be/svc-email"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	var cfg config.AppConfig

	err := config.InitConfig(&cfg)
	if err != nil {
		return err
	}
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}

	err = email.Root(s.Waiter().Context(), s)
	if err != nil {
		return err
	}
	defer fmt.Println("module shutdown")

	s.Waiter().Add(
		s.WaitForWeb,
		// s.WaitForGRPC,
	)

	return s.Waiter().Wait()

}
