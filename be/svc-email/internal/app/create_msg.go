package app

import (
	"context"
	"fmt"
	htmltemplate "html/template"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/email"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
	tmplpbv1 "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func (a *Application) createMsg(ctx context.Context, p *emailpbv1.EmailMessage) (*email.Msg, error) {

	var err error
	var msg *email.Msg

	t := p.GetBody()
	if t == nil {
		return nil, &UnsupportedBodyTypeError{
			t: "body is nil",
		}
	}
	switch v := t.(type) {
	case *emailpbv1.EmailMessage_RawBody:
		fmt.Printf("\t%#v \n", v)
		msg, err = fromRawBody(p, v)
	case *emailpbv1.EmailMessage_TemplateBody:
		fmt.Printf("\t%#v \n", v)
		msg, err = a.fromTemplate(ctx, p, v)
	default:
		return nil, &UnsupportedBodyTypeError{
			t: fmt.Sprintf("%T", v),
		}
	}
	if err != nil {
		return nil, err
	}

	return msg, err

}

func fromRawBody(cmd *emailpbv1.EmailMessage, body *emailpbv1.EmailMessage_RawBody) (*email.Msg, error) {
	msg, err := email.CreateMsgFromString(cmd.FromAccount.Email, cmd.CreateToList(), body.RawBody.Subject)
	if err != nil {
		return nil, err
	}

	msg.AddHtmlBodyString(body.RawBody.Body)

	return msg, nil

}

func (a *Application) fromTemplate(ctx context.Context, cmd *emailpbv1.EmailMessage, body *emailpbv1.EmailMessage_TemplateBody) (*email.Msg, error) {

	value, err := a.tmplFetcher.Fetch(ctx, body.TemplateBody.TemplateId)
	if err != nil {
		return nil, err
	}

	tmpl := &tmplpbv1.Template{}
	err = proto.Unmarshal(value.Value, tmpl)
	if err != nil {
		return nil, &TemplateError{err: err}
	}

	msg, err := email.CreateMsgFromString(cmd.FromAccount.Email, cmd.CreateToList(), tmpl.Name)
	if err != nil {
		return nil, err
	}

	tmplParser, err := htmltemplate.New("template_data").
		Parse(tmpl.Body)
	if err != nil {
		return nil, err
	}

	msg.AddHtmlBodyWriter(func(w io.Writer) error {
		// TODO need to convert proto data to whichever
		return tmplParser.Execute(w, body.TemplateBody.Data)
	})

	// msg.AddHtmlBodyString(tmpl.Body)

	return msg, nil

}
