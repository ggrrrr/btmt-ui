package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	NatsPublisher struct {
		conn           *NatsConnection
		subject        string
		tokenGenerator token.ServiceTokenGenerator
	}
)

func NewPublisher(conn *NatsConnection, subject string, tokenGenerator token.ServiceTokenGenerator) (*NatsPublisher, error) {

	return &NatsPublisher{
		conn:           conn,
		subject:        subject,
		tokenGenerator: tokenGenerator,
	}, nil
}

func (c *NatsPublisher) Shutdown() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.shutdown()
}

func (c *NatsPublisher) publish(ctx context.Context, uniqId string, msg *nats.Msg) error {
	var err error

	ctx, span := c.producerSpan(ctx, msg)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	ack, err := c.conn.js.PublishMsg(ctx, msg, jetstream.WithMsgID(uniqId))
	if err != nil {
		return err
	}

	addSpanAttributes(span, ack)
	return nil
}

func (c *NatsPublisher) Publish(ctx context.Context, msg *PublishMsg) error {

	subject := c.subject
	if msg.SubjectKey != "" {
		subject = fmt.Sprintf("%s.%s", subject, msg.SubjectKey)
	}

	msg.msg = nats.Msg{
		Subject: subject,
		Data:    msg.Payload,
		Header:  nats.Header{},
	}

	msg.Set("Content-Type", msg.ContentType)

	token, err := c.tokenGenerator.TokenForService(ctx)
	if err != nil {
		return fmt.Errorf("unable to sign msg: %w", err)
	}

	otel.GetTextMapPropagator().Inject(ctx, msg)

	msg.msg.Header.Set(authHeaderName, token)

	return c.publish(ctx, msg.UniqId, &msg.msg)
}
