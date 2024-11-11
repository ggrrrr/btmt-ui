package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type (
	DataHandler func(ctx context.Context, subject string, data []byte)

	NatsConsumer struct {
		conn     *Nats
		stream   jetstream.Stream
		consumer jetstream.Consumer
		shutdown func()
	}
)

func NewConsumer(ctx context.Context, conn Nats, streamName string, group string) (*NatsConsumer, error) {
	stream, err := conn.js.Stream(ctx, streamName)
	if err != nil {
		return nil, err
	}

	c, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   group,
		AckPolicy: jetstream.AckAllPolicy,
	})

	if err != nil {
		return nil, err
	}

	return &NatsConsumer{
		conn:     &conn,
		stream:   stream,
		consumer: c,
	}, nil
}

func (c *NatsConsumer) ConsumerLoop(handler DataHandler) error {
	consLoop, err := c.consumer.Consume(
		func(msg jetstream.Msg) {
			m := &jsMsg{
				msg: msg,
			}
			spanOpts := []oteltrace.SpanStartOption{
				// oteltrace.WithAttributes(logger.AttributeString("asd")),
				// oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}
			md, _ := msg.Metadata()
			if md != nil {
				attr := oteltrace.WithAttributes(
					logger.AttributeString("nats.js.domain", md.Domain),
					logger.AttributeString("nats.js.stream", md.Stream),
					logger.AttributeString("nats.js.sequence", fmt.Sprintf("%d/%d", md.Sequence.Stream, md.Sequence.Consumer)),
				)
				spanOpts = append(spanOpts, attr)
			}

			ctx := otel.GetTextMapPropagator().Extract(context.Background(), m)
			// spanCtx := oteltrace.SpanContextFromContext(ctx)
			// if spanCtx.IsValid() && spanCtx.IsRemote() {
			// 	// spanOpts = append(
			// 	// spanOpts,
			// 	// oteltrace.WithLinks(oteltrace.Link{
			// 	// SpanContext: spanCtx,
			// 	// }),
			// 	// )
			// }

			ctx, span := logger.Tracer().Start(ctx, msg.Subject(), spanOpts...)
			defer span.End()

			msg.Ack()
			handler(ctx, msg.Subject(), msg.Data())
		},
		jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
			fmt.Printf("some error, %+v %+v \n\n", consumeCtx, err)
		}),
	)
	if err != nil {
		return err
	}
	c.shutdown = consLoop.Stop
	return nil
}

func (c *NatsConsumer) Shutdown() {
	if c.shutdown != nil {
		c.shutdown()
		fmt.Printf("ConsumerLoop.Shutdown.\n")
	}
}
