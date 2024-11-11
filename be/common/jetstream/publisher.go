package jetstream

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type (
	NatsPublisher struct {
		conn    *Nats
		subject string
	}
)

func NewPublisher(conn *Nats, subject string) *NatsPublisher {

	return &NatsPublisher{
		conn:    conn,
		subject: subject,
	}
}

func (c *NatsPublisher) publish(ctx context.Context, msg *nats.Msg) error {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "jetstream.publish", nil, logger.KVString("subject", msg.Subject))
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

	return c.publish(ctx, msg.msg)
}
