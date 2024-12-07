package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ggrrrr/btmt-ui/be/common/email"
	"github.com/ggrrrr/btmt-ui/be/common/state"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

func TestSend(t *testing.T) {

	ctx := context.Background()

	tmplFetcher := state.NewMockStore()
	mockSender := &email.MockSmtpConnector{}
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
			name:     "raw",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				// mapData, err := structpb.NewStruct(map[string]any{"mapKey_1": "val 1"})
				// require.NoError(t, err)
				return &emailpbv1.EmailMessage{
					ToEmail: []*emailpbv1.ToEmail{
						{
							Name:  "to email",
							Email: "some@mail.com",
						},
					},
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "Sender",
						Email: "sender@mail.com",
					},
					Body: &emailpbv1.EmailMessage_RawBody{
						RawBody: &emailpbv1.RawBody{
							ContentType: "type",
							Subject:     "subject",
							Body:        "body",
						},
					},
				}
			},
		},
		{
			name:     "no body",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:name\nsome body val 1",
			expErr:   &UnsupportedBodyTypeError{},
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				// mapData, err := structpb.NewStruct(map[string]any{"mapKey_1": "val 1"})
				// require.NoError(t, err)
				return &emailpbv1.EmailMessage{
					ToEmail: []*emailpbv1.ToEmail{
						{
							Name:  "to email",
							Email: "some@mail.com",
						},
					},
					FromAccount: &emailpbv1.SenderAccount{
						Realm: "localhost",
						Name:  "Sender",
						Email: "sender@mail.com",
					},
				}
			},
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
