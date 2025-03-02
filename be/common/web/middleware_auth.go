package web

import (
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

func (s *Server) handlerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if s.verifier == nil {
			next.ServeHTTP(w, r)
			return
		}

		logAttr := []slog.Attr{}

		var authInfo roles.AuthInfo
		var ctx = r.Context()
		var err error
		var span = trace.SpanFromContext(ctx)
		defer func() {
			if err != nil {
				log.Log().ErrorCtx(ctx, err, "cors", logAttr...)
			} else {
				log.Log().DebugCtx(ctx, "cors", logAttr...)
			}
		}()

		userRequest := roles.FromHttpRequest(r.Header, r.Cookies(), r)

		if !userRequest.AuthData.IsZero() {
			authInfo, err = s.verifier.Verify(userRequest.AuthData)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				logAttr = append(logAttr,
					slog.String("Method", r.Method),
					slog.Any("request.device", userRequest.Device),
					slog.String("request.AuthScheme", userRequest.AuthData.AuthScheme),
				)
				SendError(ctx, w, app.UnauthenticatedError("Unauthenticated", nil))
				return
			}
			logAttr = append(logAttr, slog.Any("authInfo", authInfo))
		}

		logAttr = append(logAttr,
			slog.Any("Method", r.Method),
			slog.Any("HttpUserAgent", r.Header.Get(roles.HttpUserAgent)),
			slog.Any("remoteAddr", r.RemoteAddr),
			slog.Any("device", userRequest.Device),
			slog.Any("AuthScheme", userRequest.AuthData.AuthScheme),
			slog.Any("FullMethod", userRequest.FullMethod),
		)

		authInfo.Device = userRequest.Device
		ctx = roles.CtxWithAuthInfo(r.Context(), authInfo)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
