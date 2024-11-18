package system

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

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

type (
	httpPanicResponse struct {
		Code    int    `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}
	httpVersionReponse struct {
		BuildVersion string    `json:"build_version"`
		BuildTime    time.Time `json:"build_time"`
	}
)

func (s *System) httpHandlerRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				logger.ErrorCtx(r.Context(), fmt.Errorf("recover %+v", err)).
					// Str("stack", string(debug.Stack())).
					Msg("httpRecovery")

				fmt.Println(string(debug.Stack()))

				jsonBody, _ := json.Marshal(httpPanicResponse{
					Code:    http.StatusInternalServerError,
					Message: "internal server error",
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, err = w.Write(jsonBody)
				if err != nil {
					fmt.Printf("unable to write recover response %s", err)
				}
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
			out := "."
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

		userRequest := roles.FromHttpRequest(r.Header, r.Cookies(), r)

		if userRequest.AuthData.AuthScheme != "" {
			authInfo, err = s.verifier.Verify(userRequest.AuthData)
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
			Str("AuthScheme", userRequest.AuthData.AuthScheme).
			Str("FullMethod", userRequest.FullMethod)

		authInfo.Device = userRequest.Device
		ctx := roles.CtxWithAuthInfo(r.Context(), authInfo)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}

func (s *System) httpHandlerVersion(endpoint string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if (r.Method == "GET") && strings.EqualFold(r.URL.Path, endpoint) {
				ver := httpVersionReponse{
					BuildVersion: s.buildVersion,
					BuildTime:    s.buildTime,
				}
				web.SendPayload(r.Context(), w, "ok", ver)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}

func (s *System) initMux() {
	s.mux = chi.NewRouter()
	s.mux.NotFound(web.MethodNotFoundHandler)
	s.mux.MethodNotAllowed(web.MethodNotAllowedHandler)

	s.mux.Use(middleware.URLFormat)
	s.mux.Use(s.httpHandlerCORS)
	s.mux.Use(middleware.Logger)
	// s.mux.Use(middleware.Recoverer)
	s.mux.Use(s.httpHandlerRecoverer)
	s.mux.Use(s.httpHandlerAuth)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(middleware.Heartbeat("/liveness"))
	s.mux.Use(s.httpHandlerVersion("/_version"))

	// TODO create grpc gateway as a service
	// s.gateway = runtime.NewServeMux()

	// or use https://github.com/Ecostack/otelchi/blob/master/middleware.go
	if s.cfg.Otel.Enabled {
		s.mux.Use(otelchi.Middleware("rest"))
	}

}

func (s *System) WaitForWeb(ctx context.Context) error {
	addr := s.cfg.Rest.Address
	webServer := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	// s.mux.Mount("/v1", s.gateway)
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
