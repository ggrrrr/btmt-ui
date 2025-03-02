package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
)

func (s *Server) handlerLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// r = r.WithContext(r.Context())
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		// buf := newLimitBuffer(512)
		// ww.Tee(buf)
		// t1 := time.Now()
		defer func() {
			// var respBody []byte
			// if ww.Status() >= 400 {
			// respBody, _ = io.ReadAll(buf)
			// }
			// TODO
			logAttr := []slog.Attr{
				slog.Int("http_status", ww.Status()),
				slog.Int("size", ww.BytesWritten()),
				slog.String("TimeDiff", "ASDASDASDASDASDASDASDASD"),
				// TimeDiff("ms", time.Now(), t1),
				slog.String("path", r.URL.Path),
				slog.String("from", r.RemoteAddr),
				slog.String("http_method", r.Method),
				slog.String("http_proto", r.Proto),
			}
			span := trace.SpanFromContext(r.Context())
			if span.SpanContext().TraceID().IsValid() {
				fmt.Println("TODOTODOTODOTODOTODOTODOTODOTODO")
			}
			log.Log().Info("web.server", logAttr...)
		}()
		next.ServeHTTP(ww, r)

	})
}
