package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

const (
	spanAttributesConsumerId string = "nats.js.consumer.id"
	spanAttributesDomain     string = "nats.js.domain"
	spanAttributesDuplicate  string = "nats.js.duplicate"
	spanAttributesStream     string = "nats.js.stream"
	spanAttributesConsumer   string = "nats.js.consumer"
	spanAttributesTimestamp  string = "nats.js.timestamp"
	spanAttributesSequence   string = "nats.js.sequence"
)

func (c *NatsConsumer) consumerSpan(ctx context.Context, msg jetstream.Msg) (context.Context, oteltrace.Span) {
	spanOpts := []oteltrace.SpanStartOption{
		oteltrace.WithSpanKind(oteltrace.SpanKindConsumer),
	}

	attributes := make([]attribute.KeyValue, 1, 10)
	// We dont care if metadata fail
	attributes[0] = attribute.String("nats.js.consumer.id", c.consumerId.String())

	md, _ := msg.Metadata()
	if md != nil {
		attributes = append(attributes,
			attribute.String("nats.js.domain", md.Domain),
			attribute.String("nats.js.stream", md.Stream),
			attribute.String("nats.js.consumer", md.Consumer),
			attribute.String("nats.js.timestamp", md.Timestamp.String()),
			attribute.String("nats.js.sequence.stream", fmt.Sprintf("%d", md.Sequence.Stream)),
			attribute.String("nats.js.sequence.consumer", fmt.Sprintf("%d", md.Sequence.Consumer)),
		)
	}
	spanOpts = append(spanOpts, oteltrace.WithAttributes(attributes...))
	return logger.Tracer().Start(ctx, msg.Subject(), spanOpts...)
}

func (c *NatsConsumer) producerSpan(ctx context.Context, msg *nats.Msg) (context.Context, oteltrace.Span) {
	spanOpts := []oteltrace.SpanStartOption{
		oteltrace.WithSpanKind(oteltrace.SpanKindProducer),
	}

	attributes := make([]attribute.KeyValue, 0, 0)

	// attributes = append(attributes,
	// attribute.String("", msg.Subject),
	// )

	spanOpts = append(spanOpts, oteltrace.WithAttributes(attributes...))
	return logger.Tracer().Start(ctx, msg.Subject, spanOpts...)
}

func addSpanAttributes(span oteltrace.Span, ack *jetstream.PubAck) {
	span.SetAttributes(
		attribute.String(spanAttributesDomain, ack.Domain),
		attribute.String(spanAttributesStream, ack.Stream),
		attribute.Bool(spanAttributesDuplicate, ack.Duplicate),
		attribute.String(spanAttributesSequence, fmt.Sprintf("%v", ack.Sequence)),
	)
}
