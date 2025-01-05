package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/ggrrrr/btmt-ui/be/common/email"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	templv1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
	tmplpbv1 "github.com/ggrrrr/btmt-ui/be/svc-tmpl/tmplpb/v1"
)

func storeData(t *testing.T, store *state.MockStore, d *tmplpbv1.Template) {
	bytes, err := proto.Marshal(d)
	require.NoError(t, err)
	data := state.EntityState{
		Revision: 1,
		Key:      d.Id,
		Value:    bytes,
	}

	store.On("Fetch", data.Key).Return(data, nil)
}

func TestSend(t *testing.T) {

	ctx := context.Background()

	tmplFetcher := &state.MockStore{}
	mockSender := &email.MockSmtpConnector{
		Sender: &email.MockSmtpSender{},
	}
	testApp := &Application{
		connector:   mockSender,
		tmplFetcher: tmplFetcher,
	}

	tests := []struct {
		name     string
		from     func(t *testing.T) *emailpbv1.EmailMessage
		expected string
		expErr   error
	}{
		{
			name:     "ok from templ",
			expected: "to:to@email.com\nfrom:from@email.com\nheader:From:from@email.com\nheader:To:to@email.com\nheader:Subject:name\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				storeData(t, tmplFetcher, &tmplpbv1.Template{
					Id:          "template_id_1",
					ContentType: "type",
					Name:        "name",
					Body:        "body",
				})
				// tmplFetcher.On("Fetch", "asdasd").Return(nil, fmt.Errorf("fetch error"))

				// defer tmplFetcher.
				return &emailpbv1.EmailMessage{
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "",
						Email: "from@email.com",
					},
					ToAddresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{Name: "", Email: "to@email.com"},
						},
					},
					Body: &emailpbv1.EmailMessage_TemplateId{
						TemplateId: "template_id_1",
					},
					Data: &templv1.Data{},
				}
			},
		},
		{
			name:     "ok from payload",
			expected: "to:to@email.com\nfrom:from@email.com\nheader:From:from@email.com\nheader:To:to@email.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				return &emailpbv1.EmailMessage{
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "",
						Email: "from@email.com",
					},
					ToAddresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{Name: "", Email: "to@email.com"},
						},
					},
					Body: &emailpbv1.EmailMessage_RawBody{
						RawBody: &emailpbv1.RawBody{
							ContentType: "",
							Subject:     "subject",
							Body:        "body",
						},
					},
					Data: &templv1.Data{},
				}
			},
		},
		{
			name:     "err from templ",
			expected: "to:to@email.com\nfrom:from@email.com\nheader:From:from@email.com\nheader:To:to@email.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				tmplFetcher.On("Fetch", "asdasd").Return(nil, fmt.Errorf("fetch error"))

				return &emailpbv1.EmailMessage{
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "",
						Email: "from@email.com",
					},
					ToAddresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{Name: "", Email: "to@email.com"},
						},
					},
					Body: &emailpbv1.EmailMessage_TemplateId{
						TemplateId: "asdasd",
					},
					Data: &templv1.Data{},
				}
			},
			expErr: fmt.Errorf("some error"),
		},
		{
			name:     "err nil body",
			expected: "to:to@email.com\nfrom:from@email.com\nheader:From:from@email.com\nheader:To:to@email.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				return &emailpbv1.EmailMessage{
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "",
						Email: "from@email.com",
					},
					ToAddresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{Name: "", Email: "to@email.com"},
						},
					},
					Body: &emailpbv1.EmailMessage_RawBody{},
					Data: &templv1.Data{},
				}
			},
			expErr: fmt.Errorf("some error"),
		},
		{
			name:     "err empty body",
			expected: "to:to@email.com\nfrom:from@email.com\nheader:From:from@email.com\nheader:To:to@email.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				return &emailpbv1.EmailMessage{
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "",
						Email: "from@email.com",
					},
					ToAddresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{Name: "", Email: "to@email.com"},
						},
					},
					Body: &emailpbv1.EmailMessage_RawBody{
						RawBody: &emailpbv1.RawBody{},
					},
					Data: &templv1.Data{},
				}
			},
			expErr: fmt.Errorf("some error"),
		},
		{
			name:     "err connect",
			expected: "to:to@email.com\nfrom:from@email.com\nheader:From:from@email.com\nheader:To:to@email.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				mockSender.ForErr = fmt.Errorf("error")
				return &emailpbv1.EmailMessage{
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "",
						Email: "from@email.com",
					},
					ToAddresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{Name: "", Email: "to@email.com"},
						},
					},
					Body: &emailpbv1.EmailMessage_RawBody{
						RawBody: &emailpbv1.RawBody{
							Subject: "asd",
							Body:    "asd",
						},
					},
					Data: &templv1.Data{},
				}
			},
			expErr: fmt.Errorf("some error"),
		},
		{
			name:     "err send",
			expected: "to:to@email.com\nfrom:from@email.com\nheader:From:from@email.com\nheader:To:to@email.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				mockSender.ForErr = nil
				mockSender.Sender.ForErr = fmt.Errorf("asd")
				return &emailpbv1.EmailMessage{
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "",
						Email: "from@email.com",
					},
					ToAddresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{Name: "", Email: "to@email.com"},
						},
					},
					Body: &emailpbv1.EmailMessage_RawBody{
						RawBody: &emailpbv1.RawBody{
							Subject: "subject",
							Body:    "body",
						},
					},
					Data: &templv1.Data{},
				}
			},
			expErr: fmt.Errorf("some error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			command := tc.from(t)
			err := testApp.SendEmail(ctx, command)
			if tc.expErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expected, mockSender.Sender.LastMail)
			} else {
				require.ErrorAs(t, err, &tc.expErr)
				fmt.Printf("%s %v %#v \n\n\n", tc.name, err, err)
			}

		})
	}

}
