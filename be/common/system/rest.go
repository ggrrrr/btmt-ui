package system

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func (s *System) httpRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				logger.ErrorCtx(r.Context(), fmt.Errorf("recover %+v", err)).
					// Str("stack", string(debug.Stack())).
					Msg("httpRecovery")

				fmt.Println(string(debug.Stack()))

				jsonBody, _ := json.Marshal(struct {
					Code    int    `json:"code,omitempty"`
					Message string `json:"message,omitempty"`
					Error   string `json:"error,omitempty"`
				}{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonBody)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *System) httpHandlerCORS(next http.Handler) http.Handler {
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
		next.ServeHTTP(w, r)
	})
}

func (s *System) httpHandlerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var infoLog *zerolog.Event
		if s.cfg.Jwt.UseMock == "" {
			infoLog = logger.Debug()
		} else {
			infoLog = logger.Warn().Str("VERIFICATION", "mock")
		}
		defer infoLog.Msg("httpMiddleware")

		var authInfo roles.AuthInfo
		var err error

		userRequest := roles.FromHttp(r.Header, r.Cookies(), r.RequestURI)

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
			infoLog.Any("authInfo", authInfo)
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

	s.mux.Use(middleware.URLFormat)
	s.mux.Use(s.httpHandlerCORS)
	s.mux.Use(middleware.Logger)
	s.mux.Use(middleware.Recoverer)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(s.httpHandlerAuth)
	s.mux.Use(middleware.Heartbeat("/liveness"))

	if s.cfg.Otel.Enabled {
		s.mux.Use(otelchi.Middleware("rest"))
	}

	// TODO create grpc gateway as a service
	// s.gateway = runtime.NewServeMux()

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
