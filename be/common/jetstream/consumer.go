package jetstream

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	ConsumedMsg struct {
		Subject     string
		ContentType string
		Payload     []byte
	}

	DataHandlerFunc func(ctx context.Context, subject string, md msgbus.Metadata, data []byte) error

	NatsConsumer struct {
		tracer     tracer.OTelTracer
		consumerId uuid.UUID
		conn       *NatsConnection
		stream     jetstream.Stream
		consumer   jetstream.Consumer
		shutdown   func()
	}

	Consumer interface {
		Consume(ctx context.Context, handler DataHandlerFunc) error
		Shutdown() error
	}
)

var _ (Consumer) = (*NatsConsumer)(nil)

func NewConsumer(ctx context.Context, conn *NatsConnection, streamName string, group string) (*NatsConsumer, error) {

	if conn.verifier == nil {
		return nil, fmt.Errorf("verifier is not set")
	}

	consumerId := uuid.New()

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
		tracer:     tracer.Tracer(otelScope),
		consumerId: consumerId,
		conn:       conn,
		stream:     stream,
		consumer:   c,
	}, nil
}

func (c *NatsConsumer) Consume(ctx context.Context, handler DataHandlerFunc) error {
	log.Log().Info("ConsumerLoop.started",
		slog.String("consumer.id", c.consumerId.String()),
	)

	// c.consumer.Qon
	consLoop, err := c.consumer.Consume(
		func(msg jetstream.Msg) {
			var (
				err      error
				span     tracer.OTelSpan
				authInfo roles.AuthInfo
			)

			intMsg := &jsMsg{
				msg: msg,
			}

			ctx := otel.GetTextMapPropagator().Extract(ctx, intMsg)
			ctx, span = c.consumerSpan(ctx, msg)
			defer func() {
				span.End(err)
			}()

			authInfo, err = c.conn.verifier.Verify(app.AuthData{
				// We don`t set Auth Schema in our publisher since this is internal svc only channel
				AuthScheme: roles.AuthSchemeBearer,
				AuthToken:  intMsg.Get(authHeaderName),
			})
			if err != nil {
				log.Log().ErrorCtx(ctx, err, "verifier.Verify",
					slog.String("consumer.id", c.consumerId.String()),
				)

				return
			}
			ctx = roles.CtxWithAuthInfo(ctx, authInfo)

			md := createMessageMD(msg)

			err = msg.Ack()
			if err != nil {
				log.Log().ErrorCtx(ctx, err, "msg.Ack failed",
					slog.String("consumer.id", c.consumerId.String()),
				)
			}

			err = handler(ctx, msg.Subject(), md, msg.Data())
			if err != nil {
				// TODO now what!?
				// write to DLQ (Dead Letter Queue)
				log.Log().ErrorCtx(ctx, err, "handler",
					slog.String("consumer.id", c.consumerId.String()),
				)

			}
		},
		jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
			log.Log().ErrorCtx(ctx, err, "handler.ConsumeErrHandler",
				slog.String("consumer.id", c.consumerId.String()),
			)
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
		log.Log().Info("Shutdown",
			slog.String("consumerId", c.consumerId.String()))
		c.shutdown()
	}

	if c.conn == nil {
		return nil
	}
	return c.conn.shutdown()
}
