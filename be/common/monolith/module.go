package monolith

import (
	"context"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
)

type Monolith struct {
	Config config.AppConfig
	// DB() *sql.DB
	Mux    *chi.Mux
	RPC    *grpc.Server
	Waiter waiter.Waiter
}

type Module interface {
	Startup(context.Context, Monolith) error
}
