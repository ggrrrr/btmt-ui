package jetstream

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
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

func (c *NatsConsumer) consumerSpan(ctx context.Context, msg jetstream.Msg) (context.Context, tracer.OTelSpan) {
	logAttr := make([]slog.Attr, 0, 10)
	// We dont care if metadata fail

	logAttr = append(logAttr,
		slog.String("nats.js.consumer.id", c.consumerId.String()),
		slog.String("consumer.id", c.consumerId.String()),
		slog.Any("mg", msg.Headers()),
	)

	defer func() {
		log.Log().DebugCtx(ctx, "Consume.Headers", logAttr...)
	}()

	md, _ := msg.Metadata()
	if md != nil {
		// logAttr = append(logAttr,
		// 	slog.Any("mg", md),
		// )
		logAttr = append(logAttr,
			slog.String("nats.js.domain", md.Domain),
			slog.String("nats.js.stream", md.Stream),
			slog.String("nats.js.consumer", md.Consumer),
			slog.String("nats.js.timestamp", md.Timestamp.String()),
			slog.String("nats.js.sequence.stream", fmt.Sprintf("%d", md.Sequence.Stream)),
			slog.String("nats.js.sequence.consumer", fmt.Sprintf("%d", md.Sequence.Consumer)),
		)
	}
	// spanOpts = append(spanOpts, oteltrace.WithAttributes(attributes...))
	return c.tracer.SpanWithAttributes(ctx, msg.Subject(), logAttr...)
}

func (c *NatsPublisher) producerSpan(ctx context.Context, msg *nats.Msg) (context.Context, tracer.OTelSpan) {

	logAttr := make([]slog.Attr, 0, 10)
	logAttr = append(logAttr,
		slog.String("subject", msg.Subject),
	)

	return c.tracer.SpanWithAttributes(ctx, msg.Subject, logAttr...)
}

func addSpanAttributes(span tracer.OTelSpan, ack *jetstream.PubAck) {
	span.SetAttributes(
		slog.String(spanAttributesDomain, ack.Domain),
		slog.String(spanAttributesStream, ack.Stream),
		slog.Bool(spanAttributesDuplicate, ack.Duplicate),
		slog.String(spanAttributesSequence, fmt.Sprintf("%v", ack.Sequence)),
	)
}
