package system

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/riandyrn/otelchi"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/web"
)

// https://github.com/Ecostack/otelchi/blob/master/middleware.go
func (s *System) httpMiddlewareOtel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		ctx, span := logger.Tracer().Start(
			ctx, r.RequestURI,
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("rest", "", r)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		)
		defer span.End()
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *System) httpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Authorization")

		if r.Method == http.MethodOptions {
			out := "ok"
			w.WriteHeader(200)
			_, err := w.Write([]byte(out))
			if err != nil {
				logger.ErrorCtx(r.Context(), err)
			}
			return
		}

		var infoLog *zerolog.Event
		if s.cfg.Jwt.UseMock == "" {
			infoLog = logger.Debug()
		} else {
			infoLog = logger.Warn().Str("VERIFICATION", "mock")
		}
		defer infoLog.Msg("httpMiddleware")

		var authInfo roles.AuthInfo
		var err error

		userRequest := roles.FromHttpMetadata(r.Header, r.RequestURI)

		if userRequest.Authorization.AuthScheme != "" {
			authInfo, err = s.verifier.Verify(userRequest.Authorization)
			if err != nil {
				infoLog.
					Any("Method", r.Method).
					Any("request", userRequest).
					Err(err)
				web.SendError(r.Context(), w, app.UnauthenticatedError(err.Error(), nil))
				return
			}
			infoLog.Str("user", authInfo.User)
		}
		infoLog.
			Any("Method", r.Method).
			Any("HttpUserAgent", r.Header.Get(roles.HttpUserAgent)).
			Any("remoteAddr", r.RemoteAddr).
			Any("device", userRequest.Device).
			Str("AuthScheme", userRequest.Authorization.AuthScheme).
			Str("FullMethod", userRequest.FullMethod)

		authInfo.Device = userRequest.Device
		ctx := roles.CtxWithAuthInfo(r.Context(), authInfo)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}

func (s *System) initMux() {
	s.mux = chi.NewRouter()
	s.mux.NotFound(web.MethodNotFoundHandler)
	s.mux.MethodNotAllowed(web.MethodNotAllowedHandler)
	// s.mux.NotFound(func(w http.ResponseWriter, r *http.Request) {
	// 	// chkMw2 := r.Context().Value(ctxKey{"mw2"}).(string)
	// 	w.WriteHeader(404)
	// 	w.Write([]byte(fmt.Sprintf("sub 404 %s", "adsasd")))
	// })
	s.mux.Use(middleware.Heartbeat("/liveness"))
	s.mux.Use(middleware.Logger)
	s.mux.Use(s.httpMiddleware)

	// if s.cfg.Otel.Enabled {
	// 	// s.mux.Use(otelchi.Middleware("my-server", otelchi.WithChiRoutes(s.mux)))
	// 	s.mux.Use(s.httpMiddlewareOtel)
	// }
	if s.cfg.Otel.Enabled {
		s.mux.Use(otelchi.Middleware("rest"))
		// s.mux.Use(s.httpMiddlewareOtel)
	}

	s.gateway = runtime.NewServeMux()

	// s.mux.Get("/root", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("root."))
	// })
	// s.mux.Method("GET", "/metrics", promhttp.Handler())
}

func (s *System) WaitForWeb(ctx context.Context) error {
	addr := s.cfg.Rest.Address
	webServer := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	s.mux.Mount("/v1", s.gateway)
	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		logger.Info().Str("address", addr).Msg("rest started")
		defer fmt.Println("web server shutdown")
		if err := webServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		logger.Info().Str("address", addr).Msg("rest server to be shutdown")
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := webServer.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	})
	return group.Wait()
}
