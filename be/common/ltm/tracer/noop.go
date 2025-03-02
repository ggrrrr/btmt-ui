package tracer

import (
	"context"
	"log/slog"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

type (
	noopTracer struct{}
	noopSpan   struct{}
)

var _ (OTelTracer) = (*noopTracer)(nil)

func (n *noopTracer) SpanWithAttributes(ctx context.Context, name string, attr ...slog.Attr) (context.Context, OTelSpan) {
	return ctx, noopSpan{}
}

func (noopTracer) Span(ctx context.Context, _ string) (context.Context, OTelSpan) {
	return ctx, noopSpan{}
}

func (noopTracer) SpanWithData(ctx context.Context, _ string, _ td.TraceDataExtractor) (context.Context, OTelSpan) {
	return ctx, noopSpan{}
}

var _ (OTelSpan) = (*noopSpan)(nil)

func (noopSpan) End(error) {}

func (n noopSpan) SetAttributes(attr ...slog.Attr) {}
