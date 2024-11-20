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
	"github.com/ggrrrr/btmt-ui/be/common/token"
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
		conn       *natsConn
		stream     jetstream.Stream
		consumer   jetstream.Consumer
		verifier   token.Verifier
		shutdown   func()
	}
)

func NewConsumer(ctx context.Context, url string, streamName string, group string, verifier token.Verifier) (*NatsConsumer, error) {
	consumerId := uuid.New()
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
		consumerId: consumerId,
		conn:       conn,
		stream:     stream,
		consumer:   c,
		verifier:   verifier,
	}, nil
}

func (c *NatsConsumer) ConsumerLoop(handler DataHandler) error {
	logger.Info().
		Str("consumer.id", c.consumerId.String()).
		Msg("ConsumerLoop.started")

	// c.consumer.Qon
	consLoop, err := c.consumer.Consume(
		func(msg jetstream.Msg) {
			var err error

			intMsg := &jsMsg{
				msg: msg,
			}

			// TODO try to pass context from executor
			ctx := otel.GetTextMapPropagator().Extract(context.Background(), intMsg)
			ctx, span := c.consumerSpan(ctx, msg)
			defer func() {
				if err != nil {
					span.RecordError(err)
					span.SetStatus(codes.Error, err.Error())
				}
				span.End()
			}()

			authInfo, err := c.verifier.Verify(app.AuthData{
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
