package msgbus

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

type (
	emailSender interface {
		SendEmail(ctx context.Context, msg *emailpbv1.EmailMessage) error
	}

	server struct {
		app emailSender
	}
)

func Start(ctx context.Context, app emailSender, consumer msgbus.MessageConsumer[*emailpbv1.SendEmail]) error {

	s := &server{
		app: app,
	}

	consumer.Consumer(ctx, s.handlerSendEmail)

	return nil
}

func (s *server) handlerSendEmail(ctx context.Context, topic string, md msgbus.Metadata, sendEmail *emailpbv1.SendEmail) error {
	var err error
	ctx, span := logger.SpanWithAttributes(ctx, "handler", sendEmail, logger.TraceKVString("topic", topic))
	defer func() {
		span.End(err)
	}()

	err = s.app.SendEmail(ctx, sendEmail.Message)
	if err != nil {
		logger.ErrorCtx(ctx, err).Msg("handlerSendEmail")
	}
	logger.InfoCtx(ctx).Str("topic", topic).Msg("handler")
	return nil
}
