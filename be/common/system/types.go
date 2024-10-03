package system

import (
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/token"
	"github.com/ggrrrr/btmt-ui/be/common/waiter"
)

type System struct {
	cfg      config.AppConfig
	mux      *chi.Mux
	gateway  *runtime.ServeMux
	waiter   waiter.Waiter
	grpc     *grpc.Server
	aws      *session.Session
	verifier token.Verifier
	db       *sql.DB
}

type Service interface {
	Config() config.AppConfig
	Mux() *chi.Mux
	Waiter() waiter.Waiter
	RPC() *grpc.Server

	DB() *sql.DB
	// JS() nats.JetStreamContext
	// Mux() *chi.Mux
	// Logger() zerolog.Logger
}

type Module interface {
	Startup(context.Context, Service) error
	Name() string
}

// var _ interface = (*struct)(nil) // Compile error on missing methods
var _ Service = (*System)(nil)

func (s *System) DB() *sql.DB {
	return s.db
}

func (s *System) Config() config.AppConfig {
	return s.cfg
}

func (s *System) RPC() *grpc.Server {
	return s.grpc
}

func (s *System) Mux() *chi.Mux {
	return s.mux
}

func (s *System) Gateway() *runtime.ServeMux {
	return s.gateway
}

func (s *System) Waiter() waiter.Waiter {
	return s.waiter
}

func (s *System) Aws() *session.Session {
	return s.aws
}
