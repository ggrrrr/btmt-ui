package server

import (
	"fmt"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/system"
)

func Server() error {
	var err error
	var cfg system.Config

	config.MustParse(&cfg)
	s, err := system.NewSystem(cfg)
	if err != nil {
		return err
	}
	// err = events.Root(s.Waiter().Context(), s)
	// if err != nil {
	// 	return err
	// }

	defer fmt.Println("events module shutdown")

	s.Waiter().Add(
	// s.WaitForWeb,
	// s.WaitForGRPC,
	)

	return s.Waiter().Wait()
}
