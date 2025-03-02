package log

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

type AppLog struct {
	callerPathLevel int
	logger          *slog.Logger
	level           slog.Level
}

func (l *AppLog) IsTrace() bool {
	return l.level > slog.LevelDebug
}

func (l *AppLog) Error(err error, msg string, a ...slog.Attr) {
	l.log(slog.LevelError, err, msg, a...)
}

func (l *AppLog) ErrorCtx(ctx context.Context, err error, msg string, a ...slog.Attr) {
	l.logCtx(ctx, slog.LevelError, err, msg, a...)
}

func (l *AppLog) WarnCtx(ctx context.Context, err error, msg string, a ...slog.Attr) {
	l.logCtx(ctx, slog.LevelWarn, err, msg, a...)
}

func (l *AppLog) Warn(err error, msg string, a ...slog.Attr) {
	l.log(slog.LevelWarn, err, msg, a...)
}

func (l *AppLog) InfoCtx(ctx context.Context, msg string, a ...slog.Attr) {
	l.logCtx(ctx, slog.LevelInfo, nil, msg, a...)
}

func (l *AppLog) Info(msg string, a ...slog.Attr) {
	l.log(slog.LevelInfo, nil, msg, a...)
}

func (l *AppLog) DebugCtx(ctx context.Context, msg string, a ...slog.Attr) {
	appLogger.logCtx(ctx, slog.LevelDebug, nil, msg, a...)
}

func (l *AppLog) Debug(msg string, a ...slog.Attr) {
	appLogger.log(slog.LevelDebug, nil, msg, a...)
}

func findCaller(pathLevel int) string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown_error"
	}
	fileParts := strings.Split(file, "/")
	fromI := len(fileParts) - pathLevel
	toI := len(fileParts)
	outFile := strings.Join(fileParts[fromI:toI], "/")

	return fmt.Sprintf("%s:%d", outFile, line)
}

func (l *AppLog) logCtx(ctx context.Context, level slog.Level, err error, msg string, a ...slog.Attr) {
	addLen := 2

	td := td.Extract(ctx)
	attr := make([]slog.Attr, 0, len(a)+td.Len()+addLen)

	if l.callerPathLevel > 0 {
		attr = append(attr, slog.String("go.source.file", findCaller(l.callerPathLevel)))
	}

	attr = append(attr, td.Attr()...)
	attr = append(attr, a...)
	if err != nil {
		attr = append(attr, slog.Any("error", err))
	}
	appLogger.logger.LogAttrs(ctx, level, msg, attr...)
}

func (l *AppLog) log(level slog.Level, err error, msg string, a ...slog.Attr) {
	if err == nil {
		l.logger.LogAttrs(context.Background(), level, msg, a...)
	}
	attr := make([]slog.Attr, 0, len(a)+2)
	if l.callerPathLevel > 0 {
		attr = append(attr, slog.String("go.source.file", findCaller(l.callerPathLevel)))
	}
	attr = append(attr, slog.Any("error", err))
	l.logger.LogAttrs(context.Background(), level, msg, attr...)
}
