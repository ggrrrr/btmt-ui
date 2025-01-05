package app

import (
	"context"
	"fmt"
	"io"

	"github.com/ggrrrr/btmt-ui/be/common/email"
	"github.com/ggrrrr/btmt-ui/be/common/templ"
	templv1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

type msgData struct {
	fromAddress *emailpbv1.SenderAccount
	addresses   *emailpbv1.ToAddresses
	subject     string
	body        string
	data        *templv1.Data
}

func (a *Application) createMsg(
	ctx context.Context,
	body msgData,
) (*email.Msg, error) {
	var err error
	var msg *email.Msg

	if body.fromAddress == nil {
		return nil, fmt.Errorf("from address is nil")
	}

	if body.addresses == nil {
		return nil, fmt.Errorf("addresses addresses is nil")
	}

	if body.body == "" {
		return nil, fmt.Errorf("body is empty")
	}

	msg, err = email.CreateMsgFromString(body.fromAddress.Email, emailpbv1.CreateList(body.addresses.ToEmail))
	if err != nil {
		return nil, fmt.Errorf("createMsg %w", err)
	}
	err = msg.SetSubject(body.subject)
	if err != nil {
		return nil, fmt.Errorf("createMsg.SetSubject %w", err)
	}

	imgRender := createImageRender(ctx, body.fromAddress.Realm, a)

	htmlTempl, err := templ.NewHtml(body.body, templ.WithRenderImg(imgRender.renderImg))
	if err != nil {
		return nil, err
	}

	msg.AddHtmlBodyWriter(func(w io.Writer) error {
		return htmlTempl.Execute(w, body.data)
	})

	for k := range imgRender.images {
		msg.AddAttachment(k, func(w io.Writer) error {
			_, err := io.Copy(w, imgRender.images[k].ReadCloser)
			return err
		})
	}

	return msg, nil
}
