package system

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/ggrrrr/btmt-ui/be/common/web"
)

func (s *System) httpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if !strings.HasPrefix(r.RequestURI, "/rest") {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Authorization")

		if r.Method == http.MethodOptions {
			// logger.Info().
			// Any("headers", r.Header).
			// Msg("OPTIONS")
			// infoLog.Any("headers-r", r.Header.Get("Access-Control-Request-Headers"))
			// origin := r.Header.Get("Origin")
			// if len(origin) > 0 {
			// }
			// allowedHeaders := "Accept,Content-Type,content-type,Content-Length,Accept-Encoding,Authorization,X-Authorization,X-CSRF-Token,Origin,X-Requested-With"
			// w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
			// w.Header().Set("Access-Control-Allow-Origin", "*")
			// w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")

			// w.Header().Set("Access-Control-Allow-Origin", "*")
			// w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
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
				web.SendError(w, app.UnauthenticatedError(err.Error(), nil))
				return
				// return req, status.Error(codes.Unauthenticated, err.Error())
				// next.ServeHTTP(w, r)
			}
			infoLog.Str("user", authInfo.User)
		}
		infoLog.
			Any("Method", r.Method).
			Any("asd", r.Header.Get(roles.HttpUserAgent)).
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
