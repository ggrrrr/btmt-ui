package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/ggrrrr/btmt-ui/be/common/state"
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

func TestCreateMsg(t *testing.T) {
	ctx := context.Background()

	tmplFetcher := &state.MockStore{}

	testApp := &Application{
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
					ToEmail: []*emailpbv1.EmailAddr{
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
			name:     "tmpl",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:name\nsome body val 1",
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				mapData, err := structpb.NewStruct(map[string]any{"mapKey_1": "val 1"})
				require.NoError(t, err)

				storeData(t, tmplFetcher, &tmplpbv1.Template{
					Id:          "template_id_1",
					ContentType: "type",
					Name:        "name",
					Body:        "some body {{ .mapKey_1 }}",
				})

				return &emailpbv1.EmailMessage{
					ToEmail: []*emailpbv1.EmailAddr{
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
					Body: &emailpbv1.EmailMessage_TemplateBody{
						TemplateBody: &emailpbv1.TemplateBody{
							TemplateId: "template_id_1",
							Data:       mapData,
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

				storeData(t, tmplFetcher, &tmplpbv1.Template{
					Id:          "asd",
					ContentType: "type",
					Name:        "name",
					Body:        "some body {{ .mapKey_1 }}",
				})

				return &emailpbv1.EmailMessage{
					ToEmail: []*emailpbv1.EmailAddr{
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
		{
			name:     "no error on fetch tmpl",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:name\nsome body val 1",
			expErr:   fmt.Errorf("asdasd"),
			from: func(t *testing.T) *emailpbv1.EmailMessage {
				mapData, err := structpb.NewStruct(map[string]any{"mapKey_1": "val 1"})
				require.NoError(t, err)
				tmplFetcher.On("Fetch", "template_id_2").Return(state.EntityState{}, fmt.Errorf("asdasd"))
				return &emailpbv1.EmailMessage{
					ToEmail: []*emailpbv1.EmailAddr{
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
					Body: &emailpbv1.EmailMessage_TemplateBody{
						TemplateBody: &emailpbv1.TemplateBody{
							TemplateId: "template_id_2",
							Data:       mapData,
						},
					},
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			command := tc.from(t)
			msg, err := testApp.createMsg(ctx, command)
			if tc.expErr == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expected, msg.DumpToText())
			} else {
				require.ErrorAs(t, err, &tc.expErr)
				fmt.Printf("%s %v %#v \n", tc.name, err, err)
				fmt.Printf("%s %v %#v \n", tc.name, &tc.expErr, &tc.expErr)
			}
		})
	}

}
