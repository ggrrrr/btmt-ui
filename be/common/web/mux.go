package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

const (
	// EndpointRest    string = ""
	EndpointVersion string = "/_version"
	EndpointHealthz string = "/_healthz"
)

type (
	httpPanicResponse struct {
		Code    int    `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	httpVersionReponse struct {
		BuildVersion string `json:"build_version"`
	}

	httpReadyResponse struct {
		Status string `json:"status"`
	}
)

func (s *Server) handlerRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rvr := recover()
			if rvr != nil {
				if rvr == http.ErrAbortHandler {
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}
				// logger.ErrorCtx(r.Context(), fmt.Errorf("recover %+v", err)).
				// 	Str("stack", string(debug.Stack())).
				// 	Msg("httpRecovery")

				// fmt.Println(string(debug.Stack()))
				fmt.Printf("\nrecover: %#v\n\n", rvr)
				logger.StackCtx(r.Context(), fmt.Errorf("panic")).Any("panic", rvr).Stack().Send()

				SendJSONSystemError(r.Context(), w, "internal panic", nil, nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.CORS.Origin != "" {
			if s.cfg.CORS.Origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", s.cfg.CORS.Origin)
			}
			if s.cfg.CORS.Headers != "" {
				w.Header().Set("Access-Control-Allow-Headers", s.cfg.CORS.Headers)
			}
		}

		if r.Method == http.MethodOptions {
			sendText(r.Context(), w, 200, ".")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if s.verifier == nil {
			next.ServeHTTP(w, r)
			return
		}

		infoLog := logger.Debug()
		var authInfo roles.AuthInfo
		var err error
		var span = trace.SpanFromContext(r.Context())
		defer func() {
			infoLog.Send()
		}()

		userRequest := roles.FromHttpRequest(r.Header, r.Cookies(), r)

		if !userRequest.AuthData.IsZero() {
			authInfo, err = s.verifier.Verify(userRequest.AuthData)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				infoLog.
					Any("Method", r.Method).
					Any("request", userRequest).
					Err(err)
				SendError(r.Context(), w, app.UnauthenticatedError("Unauthenticated", nil))
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
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerVersion(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet &&
			strings.EqualFold(r.URL.Path, EndpointVersion) {
			ver := httpVersionReponse{
				BuildVersion: s.buildVersion,
			}
			SendJSONPayload(r.Context(), w, "ok", ver)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handlerReady(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet &&
			strings.EqualFold(r.URL.Path, EndpointHealthz) {

			if s.readyFunc == nil {
				sendText(r.Context(), w, 200, "undefined")
				return
			}

			if !s.readyFunc() {
				sendText(r.Context(), w, 500, "not ready")
				return
			}

			sendText(r.Context(), w, 200, "ok")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		l := logger.Info().Ctx(r.Context())

		next.ServeHTTP(w, r)

		l.
			Str("method", r.Method).
			Str("url", r.URL.RequestURI()).
			Str("user_agent", r.UserAgent()).
			Str("RemoteAddr", r.RemoteAddr).
			Dur("elapsed_ms", time.Since(start)).
			Msg("incoming request")
	})
}

func (s *Server) initMux() {
	s.mux = chi.NewRouter()
	s.mux.NotFound(notFoundHandler)
	s.mux.MethodNotAllowed(methodNotFoundHandler)

	s.mux.Use(middleware.Logger)

	s.mux.Use(s.handlerVersion)
	s.mux.Use(s.handlerReady)
	s.mux.Use(middleware.URLFormat)
	s.mux.Use(middleware.RealIP)
	s.mux.Use(otelchi.Middleware(s.name))
	s.mux.Use(s.handlerCORS)
	// s.mux.Use(middleware.Recoverer)
	s.mux.Use(s.handlerRecoverer)
	// s.mux.Use(s.requestLogger)
	s.mux.Use(s.handlerAuth)
}
