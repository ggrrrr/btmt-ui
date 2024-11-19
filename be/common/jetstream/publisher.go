package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	NatsPublisher struct {
		conn           *natsConn
		subject        string
		tokenGenerator token.ServiceTokenGenerator
	}
)

func NewPublisher(url string, subject string, tokenGenerator token.ServiceTokenGenerator) (*NatsPublisher, error) {
	conn, err := connect(url)
	if err != nil {
		return nil, err
	}

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

func (c *NatsPublisher) publish(ctx context.Context, msg *nats.Msg) error {
	var err error

	ctx, span := c.producerSpan(ctx, msg)
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	ack, err := c.conn.js.PublishMsg(ctx, msg)
	if err != nil {
		return err
	}

	addSpanAttributes(span, ack)
	return nil
}

func (c *NatsPublisher) Publish(ctx context.Context, key string, payload []byte) error {

	token, err := c.tokenGenerator.TokenForService(ctx)
	if err != nil {
		return fmt.Errorf("unable to sign msg: %w", err)
	}

	subject := c.subject
	if key != "" {
		subject = fmt.Sprintf("%s.%s", subject, key)
	}

	msg := natMsg{
		msg: &nats.Msg{
			Subject: subject,
			Data:    payload,
			Header:  nats.Header{},
		},
	}
	otel.GetTextMapPropagator().Inject(ctx, &msg)

	// TODO add auth header
	msg.msg.Header.Set(authHeaderName, token)

	return c.publish(ctx, msg.msg)
}
