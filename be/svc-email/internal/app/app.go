package app

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/blob"
	"github.com/ggrrrr/btmt-ui/be/common/email"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
	tmplpbv1 "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

type (
	Application struct {
		connectors   map[string]email.EmailConnector
		tmplFetcher  state.StateFetcher
		blobStore    blob.Store
		imagesFolder blob.BlobId
	}

	App interface {
		SendEmail(ctx context.Context, emailMsg *emailpbv1.EmailMessage) error
	}
)

func New(connectors map[string]email.EmailConnector, fetcher state.StateFetcher) (*Application, error) {
	a := &Application{
		connectors:  connectors,
		tmplFetcher: fetcher,
	}
	return a, nil
}

func (a *Application) SendEmail(ctx context.Context, emailMsg *emailpbv1.EmailMessage) error {

	if emailMsg == nil {
		return fmt.Errorf("email is nil")
	}

	connector, ok := a.connectors[emailMsg.FromAccount.Realm]
	if !ok {
		return fmt.Errorf("connectors for %s  not found", emailMsg.FromAccount.Realm)
	}

	data := msgData{
		fromAddress: emailMsg.FromAccount,
		addresses:   emailMsg.ToAddresses,
		data:        emailMsg.Data,
	}
	switch v := emailMsg.Body.(type) {
	case *emailpbv1.EmailMessage_RawBody:
		if v.RawBody == nil {
			return fmt.Errorf("EmailMessage_RawBody is nil")
		}
		data.subject = v.RawBody.Subject
		data.body = v.RawBody.Body
	case *emailpbv1.EmailMessage_TemplateId:
		tmpl, err := a.fetchTemplate(ctx, v.TemplateId)
		if err != nil {
			return &TemplateError{
				err: err,
			}
		}
		data.subject = tmpl.Name
		data.body = tmpl.Body
	default:
		return &UnsupportedBodyTypeError{
			t: fmt.Sprintf("%T", v),
		}
	}
	email, err := a.createMsg(ctx, data)
	if err != nil {
		// TODO 400 or 500 error
		return err
	}

	smtpInst, err := connector.Connect(ctx)
	if err != nil {
		// TODO 400 or 500 error
		return err
	}
	defer func() {
		err = smtpInst.Close()
		if err != nil {
			log.Log().ErrorCtx(ctx, err, "SendEmail.smtp.Close")
		}
	}()

	err = smtpInst.Send(ctx, email)
	if err != nil {
		return err
	}

	return nil
}

func (a *Application) fetchTemplate(ctx context.Context, templateId string) (*tmplpbv1.Template, error) {
	value, err := a.tmplFetcher.Fetch(ctx, templateId)
	if err != nil {
		return nil, err
	}

	tmpl := &tmplpbv1.Template{}
	err = proto.Unmarshal(value.Value, tmpl)
	if err != nil {
		return nil, &TemplateError{err: err}
	}

	return tmpl, nil
}
