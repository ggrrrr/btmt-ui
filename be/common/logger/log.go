package logger

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func initLog() {
	format := os.Getenv("LOG_FORMAT")
	levelStr := os.Getenv("LOG_LEVEL")

	level := strToLevel(levelStr)

	switch format {
	case "json":
		log = json(level)
	case "console":
		log = console(level)
	default:
		out := zerolog.NewConsoleWriter()
		out.NoColor = true
		l := zerolog.New(out).Level(zerolog.TraceLevel)
		log = l
	}

}

func traceMap(ctx context.Context) map[string]any {
	span := trace.SpanFromContext(ctx)
	traceId := span.SpanContext().TraceID()
	authInfo := roles.AuthInfoFromCtx(ctx)
	out := map[string]any{}
	out["device"] = authInfo.Device.DeviceInfo
	out["remote"] = authInfo.Device.RemoteAddr
	out["tenant"] = authInfo.Realm
	out["user"] = authInfo.User
	if traceId.IsValid() {
		out["trace"] = traceId.String()
	}
	return out
}

func addTrace(event *zerolog.Event, ctx context.Context) *zerolog.Event {
	return event.Any("trace", traceMap(ctx))
}

func DebugCtx(ctx context.Context) *zerolog.Event {
	l := log.Debug()
	return addTrace(l, ctx)
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Info() *zerolog.Event {
	return log.Info()
}
func Warn() *zerolog.Event {
	return log.Warn()
}

func Error(err error) *zerolog.Event {
	return log.Error().Err(err)
}

func InfoCtx(ctx context.Context) *zerolog.Event {
	l := log.Info()
	return addTrace(l, ctx)
}

func WarnCtx(ctx context.Context) *zerolog.Event {
	l := log.Warn()
	return addTrace(l, ctx)
}

func ErrorCtx(ctx context.Context, err error) *zerolog.Event {
	l := log.Error().Err(err)
	return addTrace(l, ctx)
}

func IsDebug() bool {
	switch log.GetLevel() {
	case zerolog.DebugLevel:
		return true
	case zerolog.TraceLevel:
		return true
	}

	return false
}

func IsTrace() bool {
	switch log.GetLevel() {
	case zerolog.TraceLevel:
		return true
	}

	return false
}
