package eda

import (
	"context"
	"fmt"

	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

type (
	senderApp interface {
		SendEmail(ctx context.Context, msg *emailpbv1.EmailMessage) error
	}

	server struct {
		app      senderApp
		consumer jetstream.Consumer
	}
)

func NewCommandHandler(ctx context.Context, app senderApp, conn *jetstream.NatsConnection) error {

	// stream name

	streamName := "svc-email"
	group := "email-sender"

	_, err := conn.CreateStream(ctx, "svc-email", "email sender", []string{
		"test.*",
	})
	require.NoError(t, err)

	consumer, err := jetstream.NewConsumer(context.Background(), conn, streamName, group)
	if err != nil {
		return err
	}

	err = consumer.Consume(ctx, func(ctx context.Context, subject string, data []byte) {
		logger.InfoCtx(ctx).Str("subject", subject).Msg("handle")
	})

	if err != nil {
		return fmt.Errorf("consume %w", err)
	}

	return nil
}
