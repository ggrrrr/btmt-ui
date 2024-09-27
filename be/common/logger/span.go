package logger

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
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

func AttributeString(k string, val string) attribute.KeyValue {
	key := fmt.Sprintf("%s.data.%s", spanKeyPrefix, k)
	return attribute.String(key, val)
}

func AttributeBool(k string, val bool) attribute.KeyValue {
	key := fmt.Sprintf("%s.data.%s", spanKeyPrefix, k)
	return attribute.Bool(key, val)
}

func AttributeInt(k string, val int) attribute.KeyValue {
	key := fmt.Sprintf("%s.data.%s", spanKeyPrefix, k)
	return attribute.Int(key, val)
}
