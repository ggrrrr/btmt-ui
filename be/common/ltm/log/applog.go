package log

import (
	"context"
	"log/slog"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

type AppLog struct {
	logger *slog.Logger
	level  slog.Level
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

func (l *AppLog) logCtx(ctx context.Context, level slog.Level, err error, msg string, a ...slog.Attr) {
	td := td.Extract(ctx)
	attr := make([]slog.Attr, 0, len(a)+td.Len()+1)
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
	attr := make([]slog.Attr, 0, len(a)+1)
	attr = append(attr, slog.Any("error", err))
	l.logger.LogAttrs(context.Background(), level, msg, attr...)
}
