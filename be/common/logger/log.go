package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/roles"
	"github.com/rs/zerolog"
)

type Config struct {
	Level  string `env:"LOG_LEVEL"`
	Format string `env:"LOG_FORMAT"`
}

var log zerolog.Logger

func traceMap(ctx context.Context) map[string]any {
	d := roles.AuthInfoFromCtx(ctx)
	out := map[string]any{}
	out["device"] = d.Device.DeviceInfo
	out["remote"] = d.Device.RemoteAddr
	out["user"] = d.User
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

func init() {
	out := zerolog.NewConsoleWriter()
	out.NoColor = true
	l := zerolog.New(out).Level(zerolog.TraceLevel)
	log = l
}

func console(level zerolog.Level) zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}).
		Level(level).
		With().
		Timestamp().
		Caller().
		// Int("pid", os.Getpid()).
		// Str("go_version", buildInfo.GoVersion).
		Logger()
	return logger
}

func json(level zerolog.Level) zerolog.Logger {
	var logger = zerolog.New(os.Stdout).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()
	return logger
}

func strToLevel(l string) zerolog.Level {
	switch strings.ToLower(l) {
	case "info":
		return zerolog.InfoLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.DebugLevel
	}
}

func Init(cfg Config) {
	level := strToLevel(cfg.Level)
	switch strings.ToLower(cfg.Format) {
	case "json":
		log = json(level)
	default:
		log = console(level)
	}
}
