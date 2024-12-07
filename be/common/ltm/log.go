package ltm

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracedata"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	AppLogger struct {
		logger zerolog.Logger
	}
	LogCfg struct {
		Level  string `envconfig:"LOG_LEVEL"`
		Format string `envconfig:"LOG_FORMAT"`
	}
)

func NewStdLogger(moduleName string) *AppLogger {
	return newLogger(moduleName, os.Stderr)
}

func (a *AppLogger) DebugCtx(ctx context.Context) *zerolog.Event {
	event := log.Debug()
	if ctx != nil {
		addTrace(event, ctx)
	}
	return event
}

func traceDataFromCtx(ctx context.Context) map[string]any {
	span := trace.SpanFromContext(ctx)
	traceId := span.SpanContext().TraceID()
	authInfo := roles.AuthInfoFromCtx(ctx)
	out := map[string]any{}
	out["device"] = authInfo.Device.DeviceInfo
	out["remote"] = authInfo.Device.RemoteAddr
	out["tenant"] = authInfo.Realm
	out["subject"] = authInfo.Subject
	if traceId.IsValid() {
		out["trace"] = traceId.String()
	}

	td := tracedata.TraceDataFromCtx(ctx)
	for k, v := range td {
		out[k] = v.Value()
	}

	return out

}

func addTrace(event *zerolog.Event, ctx context.Context) *zerolog.Event {
	return event.Any("trace", traceDataFromCtx(ctx))
}

func newLogger(moduleName string, writer io.WriteCloser) *AppLogger {
	var cfg LogCfg
	err := envconfig.Process(moduleName, &cfg)
	if err != nil {
		fmt.Printf("error parsing cfg for:%s", moduleName)
		cfg.Level = defaultCfg.Format
		cfg.Format = defaultCfg.Format
	}
	fmt.Printf("cfg %+v\n", cfg)

	return configureLog(cfg, writer)
}

var defaultCfg LogCfg

func defaultLog() {
	err := envconfig.Process("", &defaultCfg)
	if err != nil {
		fmt.Printf("unable to parse log config")
		panic(err)
	}
}

func configureLog(cfg LogCfg, writer io.WriteCloser) *AppLogger {
	var logger zerolog.Logger

	switch strings.ToUpper(cfg.Format) {
	case "JSON":
		logger = jsonLogger(cfg, writer)
	default:
		logger = textLogger(cfg, writer)

	}

	return &AppLogger{
		logger: logger,
	}
}

func jsonLogger(cfg LogCfg, writer io.WriteCloser) zerolog.Logger {
	var logger = zerolog.New(writer).
		Level(logLevel(cfg)).
		With().
		Timestamp().
		Caller().
		Logger()
	return logger
}

func textLogger(cfg LogCfg, writer io.WriteCloser) zerolog.Logger {

	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        writer,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}).
		Level(logLevel(cfg)).
		With().
		Timestamp().
		Caller().
		// Int("pid", os.Getpid()).
		// Str("go_version", buildInfo.GoVersion).
		Logger()
	return logger
}

func logLevel(cfg LogCfg) zerolog.Level {
	switch strings.ToUpper(cfg.Level) {
	case "DEBUG", "DBG":
		return zerolog.DebugLevel
	case "INFO", "INF":
		return zerolog.InfoLevel
	case "ERROR", "ERR":
		return zerolog.ErrorLevel
	case "WARN":
		return zerolog.WarnLevel
	case "TRACE":
		return zerolog.TraceLevel
	}
	return zerolog.InfoLevel
}
