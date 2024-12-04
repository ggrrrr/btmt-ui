package jetstream

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/ggrrrr/btmt-ui/be/common/app"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/roles"
)

type (
	ConsumedMsg struct {
		Subject     string
		ContentType string
		Payload     []byte
	}

	DataHandler func(ctx context.Context, subject string, data []byte)

	NatsConsumer struct {
		consumerId uuid.UUID
		conn       *NatsConnection
		stream     jetstream.Stream
		consumer   jetstream.Consumer
		shutdown   func()
	}

	Consumer interface {
		Consume(ctx context.Context, handler DataHandler) error
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
		consumerId: consumerId,
		conn:       conn,
		stream:     stream,
		consumer:   c,
	}, nil
}

func (c *NatsConsumer) Consume(ctx context.Context, handler DataHandler) error {
	logger.Info().
		Str("consumer.id", c.consumerId.String()).
		Msg("ConsumerLoop.started")
	// TODO use ctx for shutdown

	// c.consumer.Qon
	consLoop, err := c.consumer.Consume(
		func(msg jetstream.Msg) {
			var err error

			intMsg := &jsMsg{
				msg: msg,
			}

			// TODO try to pass context from executor
			ctx := otel.GetTextMapPropagator().Extract(ctx, intMsg)
			ctx, span := c.consumerSpan(ctx, msg)
			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
				span.End()
			}()

			authInfo, err := c.conn.verifier.Verify(app.AuthData{
				// We don`t set Auth Schema in our publisher since this is internal svc only channel
				AuthScheme: roles.AuthSchemeBearer,
				AuthToken:  intMsg.Get(authHeaderName),
			})
			if err != nil {
				logger.ErrorCtx(ctx, err).Msg("verifier.Verify")
				return
			}
			ctx = roles.CtxWithAuthInfo(ctx, authInfo)

			err = msg.Ack()
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
		fmt.Printf("Consumer[%s].Shutdown.\n", c.consumerId.String())
	}

	if c.conn == nil {
		return nil
	}
	return c.conn.shutdown()
}
