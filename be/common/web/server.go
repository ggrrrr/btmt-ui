package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/sync/errgroup"

	"github.com/ggrrrr/btmt-ui/be/common/buildversion"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	ServerOptionFn func(s *Server) error
	CORS           struct {
		Origin  string `env:"CORS_ORIGIN" envDefault:"*"`
		Headers string `env:"CORS_HEADERS" envDefault:"Content-Type, Authorization, X-Authorization"`
	}

	Config struct {
		// EndpointREST    string        `env:"ENDPOINT_REST" default:"rest"`
		ListenAddr      string        `env:"LISTEN_ADDR" envDefault:":8080"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"1s"`
		CORS            CORS
	}

	Server struct {
		name         string
		cfg          Config
		verifier     token.Verifier
		mux          *chi.Mux
		buildVersion string
		server       *http.Server
		readyFunc    func() bool
	}
)

func NewServer(name string, cfg Config, opts ...ServerOptionFn) (*Server, error) {
	if name == "" {
		name = "newServer"
	}

	// if cfg.EndpointREST == "" {
	// cfg.EndpointREST = "rest"
	// }

	if cfg.ListenAddr == "" {
		return nil, fmt.Errorf("empty ListenAddr")
	}

	s := &Server{
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

	logger.Info().
		Str("web.server", s.name).
		Any("cfg", cfg).
		Send()

	s.initMux()

	s.server = &http.Server{
		Addr:    s.cfg.ListenAddr,
		Handler: s.mux,
	}

	return s, nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info().
		Str("web.server", s.name).
		Str("ListenAddr", s.cfg.ListenAddr).
		Msg("Shutdown.")
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
	logger.Info().
		Str("web.server", s.name).
		Str("ListenAddr", s.cfg.ListenAddr).
		Msg("startup")
	defer logger.Info().
		Str("web.server", s.name).
		Msg("exit")
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("ListenAndServe: %#v \n", err)
		return err
	}
	return nil
}

func (s *Server) WaitForWeb1(ctx context.Context) error {
	// s.mux.Mount("/v1", s.gateway)
	ctx, _ = context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		return s.Startup()
	})

	group.Go(func() error {
		logger.Info().
			Str("web.server", s.name).
			Str("ListenAddr", s.cfg.ListenAddr).
			Msg("waiting...")
		<-gCtx.Done()
		logger.Info().
			Str("web.server", s.name).
			Str("ListenAddr", s.cfg.ListenAddr).
			Msg("rest server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})
	return group.Wait()
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
