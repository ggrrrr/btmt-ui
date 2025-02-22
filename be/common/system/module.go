package system

import (
	"context"
)

type Module interface {
	Configure(context.Context, Service) error
	Startup(context.Context) error
	Name() string
}
