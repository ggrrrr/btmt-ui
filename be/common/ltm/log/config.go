package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type (
	Config struct {
		Level     string `env:"LOG_LEVEL"`
		Format    string `env:"LOG_FORMAT"`
		AddSource int    `env:"LOG_ADD_SOURCE"`
	}
)

var (
	cfgLock sync.Mutex = sync.Mutex{}

	appLogger *AppLog = &AppLog{
		logger: slog.Default(),
		level:  slog.LevelDebug,
	}
)

func Log() *AppLog {
	return appLogger
}

func WithAny(k string, v any) slog.Attr {
	return slog.Any(k, v)
}

func WithBool(k string, v bool) slog.Attr {
	return slog.Bool(k, v)
}

func WithInt(k string, v int) slog.Attr {
	return slog.Int(k, v)
}

func WithString(k, v string) slog.Attr {
	return slog.String(k, v)
}

func Configure(cfg Config) error {
	cfgLock.Lock()
	defer cfgLock.Unlock()

	fmt.Printf("ltm.log: %+v\n", cfg)
	appLogger = configureWithWriter(cfg, os.Stdout)

	return nil
}

func configureWithWriter(cfg Config, writer io.Writer) *AppLog {
	out := &AppLog{
		callerPathLevel: cfg.AddSource,
		level:           parseLogLevel(cfg.Level),
	}

	opts := &slog.HandlerOptions{
		// ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {},
		// AddSource: true,
		Level: out.level,
	}

	switch strings.ToLower(cfg.Format) {
	case "json":
		out.logger = slog.New(slog.NewJSONHandler(writer, opts))
	default:
		out.logger = slog.New(slog.NewTextHandler(writer, opts))
	}

	slog.SetDefault(out.logger)

	return out
}

func parseLogLevel(from string) slog.Level {
	switch strings.ToLower(from) {
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}
