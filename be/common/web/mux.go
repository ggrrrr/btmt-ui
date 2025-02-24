package web

import (
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
)

func (s *Server) initMux() {
	s.mux = chi.NewRouter()
	s.mux.NotFound(notFoundHandler)
	s.mux.MethodNotAllowed(methodNotFoundHandler)

	s.mux.Use(s.handlerVersion)
	s.mux.Use(s.handlerReady)

	// s.mux.Use(middleware.Recoverer)
	s.mux.Use(s.handlerRecoverer)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(middleware.URLFormat)
	s.mux.Use(otelchi.Middleware(s.name))
	s.mux.Use(s.handlerLogging)

	s.mux.Use(cors.Handler(cors.Options{
		// AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
		AllowedOrigins:   strings.Split(s.cfg.CORS.Origin, ","),
	}))
	s.mux.Use(s.handlerAuth)
}
