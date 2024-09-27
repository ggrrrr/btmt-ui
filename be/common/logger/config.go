package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type traceDataCtxKey struct{}

type AttributeExtractor interface {
	Extractor() map[string]string
}

func init() {
	initLog()
	initNoopOtel()
}

func console(level zerolog.Level) zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
		NoColor:    true,
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
