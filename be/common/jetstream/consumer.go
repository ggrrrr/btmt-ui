package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	DataHandler func(ctx context.Context, subject string, data []byte)

	NatsConsumer struct {
		conn     *natsConn
		stream   jetstream.Stream
		consumer jetstream.Consumer
		verifier token.Verifier
		shutdown func()
	}
)

func NewConsumer(ctx context.Context, url string, streamName string, group string) (*NatsConsumer, error) {
	conn, err := connect(url)
	if err != nil {
		return nil, err
	}

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
		conn:     conn,
		stream:   stream,
		consumer: c,
	}, nil
}

func (c *NatsConsumer) ConsumerLoop(handler DataHandler) error {
	consLoop, err := c.consumer.Consume(
		func(msg jetstream.Msg) {
			intMsg := &jsMsg{
				msg: msg,
			}

			// TODO how to decorate the root span ?
			spanOpts := []oteltrace.SpanStartOption{
				// oteltrace.WithAttributes(logger.AttributeString("asd")),
				// oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			}

			// We dont care if metadata fail
			md, _ := msg.Metadata()
			if md != nil {
				attr := oteltrace.WithAttributes(
					attribute.String("nats.js.domain", md.Domain),
					attribute.String("nats.js.stream", md.Stream),
					attribute.String("nats.js.sequence.stream", fmt.Sprintf("%d", md.Sequence.Stream)),
					attribute.String("nats.js.sequence.consumer", fmt.Sprintf("%d", md.Sequence.Consumer)),
				)
				spanOpts = append(spanOpts, attr)
			}

			// TODO try to pass context from executor
			ctx := otel.GetTextMapPropagator().Extract(context.Background(), intMsg)
			// TODO parse auth header to include AuthInfo in the context

			// authValue := intMsg.Get(authHeaderName)

			// c.verifier.Verify(roles.Authorization{
			// 	AuthScheme: "",
			// 	AuthToken:  authValue,
			// })

			ctx, span := logger.Tracer().Start(ctx, msg.Subject(), spanOpts...)
			defer span.End()

			err := msg.Ack()
			if err != nil {
				logger.ErrorCtx(ctx, err).Msg("msg.Ack failed")
			}
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

func (c *NatsConsumer) Shutdown() error {
	if c.shutdown != nil {
		c.shutdown()
		fmt.Printf("ConsumerLoop.Shutdown.\n")
	}

	if c.conn == nil {
		return nil
	}
	return c.conn.shutdown()
}
