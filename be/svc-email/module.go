package email

import (
	"context"
	"os"

	"github.com/ggrrrr/btmt-ui/be/common/config"
	"github.com/ggrrrr/btmt-ui/be/common/email"
	"github.com/ggrrrr/btmt-ui/be/common/jetstream"
	"github.com/ggrrrr/btmt-ui/be/common/msgbus"
	"github.com/ggrrrr/btmt-ui/be/common/system"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
	"github.com/ggrrrr/btmt-ui/be/svc-email/internal/app"
	"github.com/ggrrrr/btmt-ui/be/svc-email/internal/bus"
)

type (
	Config struct {
		// Web web.Config
		BrokerCfg jetstream.Config `envPrefix:"EMAIL_"`
	}

	Module struct {
		cfg      Config
		app      app.App
		system   system.Service
		consumer msgbus.MessageConsumer[*emailpbv1.SendEmail]
	}
)

var _ (system.Module) = (*Module)(nil)

func (*Module) Name() string {
	return "svc-email"
}

func (m *Module) Configure(ctx context.Context, s system.Service) error {
	var err error
	config.MustParse(&m.cfg)
	m.system = s

	senderCfgs := map[string]email.EmailConnector{
		"localhost": &email.Config{
			SMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
			SMTPAddr: os.Getenv("EMAIL_SMTP_ADDR"),
			Username: os.Getenv("EMAIL_USERNAME"),
			Password: os.Getenv("EMAIL_PASSWORD"),
			AuthType: email.AuthTypePlain,
		},
	}

	m.app, err = app.New(senderCfgs, nil)
	if err != nil {
		return err
	}

	conn, err := jetstream.Connect(m.cfg.BrokerCfg, jetstream.WithVerifier(s.Verifier()))
	if err != nil {
		return err
	}

	m.consumer, err = jetstream.NewCommandConsumer(ctx, "svc-email", conn, &emailpbv1.SendEmail{},
		func(_ string) *emailpbv1.SendEmail { return new(emailpbv1.SendEmail) },
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Module) Startup(ctx context.Context) (err error) {

	m.system.Waiter().Add(func(ctx context.Context) error {
		return bus.Start(ctx, m.app, m.consumer)
	})

	m.system.Waiter().AddCleanup(func() {
		m.consumer.Shutdown()
	})

	return nil
}
