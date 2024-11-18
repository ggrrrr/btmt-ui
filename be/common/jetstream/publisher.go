package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	NatsPublisher struct {
		conn           *natsConn
		subject        string
		tokenGenerator token.ServiceTokenGenerator
		// we need signer to add auth header to each message
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
	ctx, span := logger.SpanWithAttributes(ctx, "jetstream.publish", nil, logger.TraceKVString("subject", msg.Subject))
	defer func() {
		span.End(err)
	}()
	_, err = c.conn.js.PublishMsg(ctx, msg)
	if err != nil {
		return err
	}
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
	msg.msg.Header.Set(authHeaderName, fmt.Sprintf("%s %s", authSchemeBearer, token))

	return c.publish(ctx, msg.msg)
}
