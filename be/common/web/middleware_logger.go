package web

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
)

func (s *Server) handlerLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// r = r.WithContext(r.Context())
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		// buf := newLimitBuffer(512)
		// ww.Tee(buf)
		startTime := time.Now()
		defer func() {
			// var respBody []byte
			// if ww.Status() >= 400 {
			// respBody, _ = io.ReadAll(buf)
			// }
			// TODO

			logAttr := []slog.Attr{
				slog.Int("http_status", ww.Status()),
				slog.Int("size", ww.BytesWritten()),
				slog.Duration("duration", time.Since(startTime)),
				slog.String("path", r.URL.Path),
				slog.String("from", r.RemoteAddr),
				slog.String("http_method", r.Method),
				slog.String("http_proto", r.Proto),
			}
			log.Log().InfoCtx(r.Context(), "web.server", logAttr...)
		}()
		next.ServeHTTP(ww, r)

	})
}
