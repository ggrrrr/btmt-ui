package app

import (
	"context"

	"github.com/ggrrrr/btmt-ui/be/common/email"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

type (
	Application struct {
		connector   email.SmtpConnector
		tmplFetcher state.StateFetcher
	}
)

func New(connector email.SmtpConnector, fetcher state.StateFetcher) (*Application, error) {
	a := &Application{
		connector: connector,
	}
	return a, nil
}

func (a *Application) SendEmail(ctx context.Context, emailMsg *emailpbv1.EmailMessage) error {
	email, err := a.createMsg(ctx, emailMsg)
	if err != nil {
		// TODO 400 or 500 error
		return err
	}

	senderInst, err := a.connector.Connect(ctx)
	if err != nil {
		// TODO 400 or 500 error
		return err
	}

	err = senderInst.Send(ctx, email)
	if err != nil {
		return err
	}

	return nil
}
