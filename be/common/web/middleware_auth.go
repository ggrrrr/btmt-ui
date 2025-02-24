package web

import (
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

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
			infoLog.Msg("handlerAuth")
		}()

		userRequest := roles.FromHttpRequest(r.Header, r.Cookies(), r)

		if !userRequest.AuthData.IsZero() {
			authInfo, err = s.verifier.Verify(userRequest.AuthData)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				infoLog.
					Str("Method", r.Method).
					Any("request.device", userRequest.Device).
					Str("request.AuthScheme", userRequest.AuthData.AuthScheme).
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
