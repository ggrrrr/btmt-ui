package web

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ggrrrr/btmt-ui/be/common/buildversion"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	ServerOptionFn func(s *Server) error
	CORS           struct {
		Origin  string `env:"WEB_CORS_ORIGIN" envDefault:"*"`
		Headers string `env:"WEB_CORS_HEADERS" envDefault:"Content-Type, Authorization, X-Authorization"`
	}

	Config struct {
		// EndpointREST    string        `env:"ENDPOINT_REST" default:"rest"`
		ListenAddr      string        `env:"WEB_LISTEN_ADDR" envDefault:":8080"`
		ShutdownTimeout time.Duration `env:"WEB_SHUTDOWN_TIMEOUT" envDefault:"1s"`
		CORS            CORS
	}

	Server struct {
		listenReady  sync.WaitGroup
		name         string
		cfg          Config
		verifier     token.Verifier
		mux          *chi.Mux
		buildVersion string
		server       *http.Server
		listener     net.Listener
		readyFunc    func() bool
	}
)

func NewServer(name string, cfg Config, opts ...ServerOptionFn) (*Server, error) {
	if name == "" {
		name = "newServer"
	}

	if cfg.ListenAddr == "" {
		return nil, fmt.Errorf("empty ListenAddr")
	}

	s := &Server{
		listenReady:  sync.WaitGroup{},
		name:         name,
		cfg:          cfg,
		buildVersion: buildversion.BuildVersion(),
	}

	for _, fn := range opts {
		err := fn(s)
		if err != nil {
			return nil, err
		}
	}

	log.Log().Info("web.server",
		slog.String("name", s.name),
		slog.Any("cfg", cfg))

	s.initMux()

	s.server = &http.Server{
		Addr:    s.cfg.ListenAddr,
		Handler: s.mux,
	}
	s.listenReady.Add(1)
	return s, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Log().Info("web.server.Shutdown",
		slog.String("name", s.name))

	ctx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *Server) MountHandler(pattern string, router http.Handler) {
	s.mux.Mount(pattern, router)
}

func (s *Server) Startup() error {
	var err error
	s.listener, err = net.Listen("tcp", s.cfg.ListenAddr)
	s.listenReady.Done()
	if err != nil {
		return err
	}

	err = s.server.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func WithVerifier(v token.Verifier) ServerOptionFn {
	return func(s *Server) error {
		s.verifier = v
		return nil
	}
}

func WithReadyFunc(readyFn func() bool) ServerOptionFn {
	return func(s *Server) error {
		s.readyFunc = readyFn
		return nil
	}
}
