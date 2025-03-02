package tracer

import (
	"context"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/td"
)

type (
	realTracer struct {
		tracer trace.Tracer
	}

	realSpan struct {
		span trace.Span
	}
)

var _ (OTelTracer) = (*realTracer)(nil)

// TODO fix tracedata
// Span implements OTelTracer.
func (r *realTracer) Span(ctx context.Context, name string) (context.Context, OTelSpan) {
	kv := td.Extract(ctx)

	ctx, span := r.tracer.Start(ctx, name, trace.WithAttributes(toAttr(kv)...))

	return ctx, &realSpan{
		span: span,
	}
}

// Span implements OTelTracer.
func (r *realTracer) SpanWithData(ctx context.Context, name string, data td.TraceDataExtractor) (context.Context, OTelSpan) {
	kv := td.Extract(ctx)

	ctx, span := r.tracer.Start(ctx, name, trace.WithAttributes(toAttr(kv)...))

	return ctx, &realSpan{
		span: span,
	}
}

func (r *realTracer) SpanWithAttributes(ctx context.Context, name string, attr ...slog.Attr) (context.Context, OTelSpan) {
	kv := td.Extract(ctx)

	ctx, span := r.tracer.Start(ctx, name, trace.WithAttributes(toAttr(kv)...))

	return ctx, &realSpan{
		span: span,
	}

}

var _ (OTelSpan) = (*realSpan)(nil)

func (a *realSpan) End(err error) {
	if err != nil {
		a.span.RecordError(err)
		a.span.SetStatus(codes.Error, err.Error())
	}
	a.span.End()
}

// SetAttributes implements OTelSpan.
func (n *realSpan) SetAttributes(attr ...slog.Attr) {
	n.span.SetAttributes(toAttr1(attr...)...)
}

func toAttr(kv *td.TraceData) []attribute.KeyValue {
	if kv.Len() == 0 {
		return []attribute.KeyValue{}
	}
	out := make([]attribute.KeyValue, 0, kv.Len())
	for i := range kv.KV {
		out = append(out,
			toKeyValue(i, kv.KV[i]),
		)
	}
	return out
}

func toAttr1(attr ...slog.Attr) []attribute.KeyValue {
	out := make([]attribute.KeyValue, 0, len(attr))
	for i := range attr {
		out = append(out,
			toKeyValue(attr[i].Key, attr[i].Value),
		)
	}
	return out
}

func toKeyValue(k string, v slog.Value) attribute.KeyValue {
	switch v.Kind() {
	case slog.KindBool:
		return attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.BoolValue(v.Bool()),
		}
	case slog.KindInt64, slog.KindUint64:
		return attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.Int64Value(v.Int64()),
		}
	case slog.KindTime:
		return attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v.Time().String()),
		}
	case slog.KindString:
		return attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v.String()),
		}
	default:
		return attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(fmt.Sprintf("%v", v.Any())),
		}
	}
}
