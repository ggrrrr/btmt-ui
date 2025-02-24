package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

func (s *Server) handlerLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// r = r.WithContext(r.Context())
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		// buf := newLimitBuffer(512)
		// ww.Tee(buf)
		t1 := time.Now()
		defer func() {
			// var respBody []byte
			// if ww.Status() >= 400 {
			// respBody, _ = io.ReadAll(buf)
			// }
			logentry := logger.Info().
				Int("http_status", ww.Status()).
				Int("size", ww.BytesWritten()).
				TimeDiff("ms", time.Now(), t1).
				Str("path", r.URL.Path).
				Str("from", r.RemoteAddr).
				Str("http_method", r.Method).
				Str("http_proto", r.Proto)
			span := trace.SpanFromContext(r.Context())
			if span.SpanContext().TraceID().IsValid() {
				logentry.Str("trace.id", span.SpanContext().TraceID().String())
			}
			// Str("http_agent", r.Header.Get("User-Agent")).
			logentry.Msg("web.server")
			// entry.Write(ww.Status(), ww.BytesWritten(), ww.Header(), time.Since(t1), respBody)
		}()
		next.ServeHTTP(ww, r)

	})
}
