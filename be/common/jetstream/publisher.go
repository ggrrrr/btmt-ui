package jetstream

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	"github.com/ggrrrr/btmt-ui/be/common/token"
)

type (
	NatsPublisher struct {
		tracer         tracer.OTelTracer
		conn           *NatsConnection
		topic          string
		tokenGenerator token.ServiceTokenGenerator
	}
)

func NewPublisher(conn *NatsConnection, topic string, tokenGenerator token.ServiceTokenGenerator) (*NatsPublisher, error) {
	log.Log().Info("NewPublisher",
		slog.String("topic", topic),
	)
	return &NatsPublisher{
		tracer:         tracer.Tracer(otelScope),
		conn:           conn,
		topic:          topic,
		tokenGenerator: tokenGenerator,
	}, nil
}

// topic is main part of the subject
// if there is OrderKey in MD then will <topic>.<orderKey>
func (c *NatsPublisher) Publish(ctx context.Context, md msgbus.Metadata, payload []byte) error {

	msg := &publishMsg{
		md:      md,
		payload: payload,
	}

	subject := c.topic
	if md.Id != uuid.Nil {
		subject = fmt.Sprintf("%s.%s", c.topic, msg.md.Id.String())
	}

	msg.msg = nats.Msg{
		Subject: subject,
		Data:    msg.payload,
		Header:  nats.Header{},
	}

	token, err := c.tokenGenerator.TokenForService(ctx)
	if err != nil {
		return fmt.Errorf("unable to sign msg: %w", err)
	}

	msg.msg.Header.Set(authHeaderName, token)
	injectHeaders(md, *msg)
	otel.GetTextMapPropagator().Inject(ctx, msg)

	return c.publish(ctx, msg.md.Id.String(), &msg.msg)
}

func (c *NatsPublisher) Shutdown() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.shutdown()
}

func (c *NatsPublisher) publish(ctx context.Context, uniqId string, msg *nats.Msg) error {

	log.Log().Info("publish",
		slog.String("subject", msg.Subject),
	)

	var err error

	ctx, span := c.producerSpan(ctx, msg)
	defer func() {
		span.End(err)
	}()

	ack, err := c.conn.js.PublishMsg(ctx, msg, jetstream.WithMsgID(uniqId))
	if err != nil {
		fmt.Printf("ack: %#v err: %#v \n", ack, err)
		return fmt.Errorf("js.PublishMsg %w", err)
	}

	addSpanAttributes(span, ack)
	return nil
}
