package email

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/email"
	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
	"github.com/ggrrrr/btmt-ui/be/svc-email/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-email/internal/msgbus"
	// "github.com/ggrrrr/btmt-ui/be/svc-auth/internal/grpc"
	// "github.com/ggrrrr/btmt-ui/be/svc-auth/internal/rest"
)

type Module struct{}

func (Module) Startup(ctx context.Context, s *system.System) (err error) {
	return Root(ctx, s)
}

func Root(ctx context.Context, s *system.System) error {
	logger.Info().Msg("Root")

	var emailApp *app.Application
	var err error

	senderCfgs := map[string]email.EmailConnector{
		// "localhost": email.Config{
		// 	SMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
		// 	SMTPAddr: os.Getenv("EMAIL_SMTP_ADDR"),
		// 	Username: os.Getenv("EMAIL_USERNAME"),
		// 	Password: os.Getenv("EMAIL_PASSWORD"),
		// 	AuthType: email.AuthTypePlain,
		// },
	}

	emailApp, err = app.New(senderCfgs, nil)
	if err != nil {
		return err
	}

	cfg := jetstream.Config{
		URL: "localhost:4222",
	}

	conn, err := jetstream.Connect(cfg, jetstream.WithVerifier(s.Verifier()))
	if err != nil {
		return err
	}

	consumer, err := jetstream.NewCommandConsumer(ctx, "svc-email", conn, &emailpbv1.SendEmail{},
		func(_ string) *emailpbv1.SendEmail { return new(emailpbv1.SendEmail) },
	)
	if err != nil {
		return err
	}
	// add shutdown func
	s.Waiter().Cleanup(func() {
		consumer.Shutdown()
	})

	msgbus.Start(ctx, emailApp, consumer)

	return nil
}
