package logger

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type appSpan struct {
	span trace.Span
}

type AppSpan interface {
	End(error)
}

var _ (AppSpan) = (*appSpan)(nil)

func (a *appSpan) End(err error) {
	if err != nil {
		a.span.RecordError(err)
		a.span.SetStatus(codes.Error, err.Error())
	}
	a.span.End()
}
