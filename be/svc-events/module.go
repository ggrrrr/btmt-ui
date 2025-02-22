package events

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/system"
)

type Module struct{}

// Configure implements system.Module.
func (m *Module) Configure(context.Context, system.Service) error {
	panic("unimplemented")
}

// Name implements system.Module.
func (m *Module) Name() string {
	panic("unimplemented")
}

// Startup implements system.Module.
func (m *Module) Startup(context.Context) error {
	panic("unimplemented")
}

var _ (system.Module) = (*Module)(nil)

// func (*Module) Startup(ctx context.Context) (err error) {
// 	return nil
// }

// func (*Module) Configure(ctx context.Context, s system.Service) error {
// 	return nil
// }
