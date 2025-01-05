package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"

	templv1 "github.com/ggrrrr/btmt-ui/be/common/templ/v1"
	emailpbv1 "github.com/ggrrrr/btmt-ui/be/svc-email/emailpb/v1"
)

func TestCreateMsg(t *testing.T) {
	ctx := context.Background()

	testApp := &Application{}

	tests := []struct {
		name     string
		from     func(t *testing.T) msgData
		expected string
		expErr   error
	}{
		{
			name:     "ok and data",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody val 1",
			from: func(t *testing.T) msgData {

				mapData, err := structpb.NewStruct(map[string]any{"mapKey_1": "val 1"})
				require.NoError(t, err)

				return msgData{
					fromAddress: &emailpbv1.SenderAccount{
						Realm: "",
						Name:  "",
						Email: "sender@mail.com",
					},
					addresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{
								Name:  "",
								Email: "some@mail.com",
							},
						},
					},
					subject: "subject",
					body:    "body {{ .Items.data.mapKey_1 }}",
					data: &templv1.Data{
						Items: map[string]*structpb.Struct{
							"data": mapData,
						},
					},
				}
			},
			expErr: nil,
		},
		{
			name:     "ok no data",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) msgData {
				return msgData{
					fromAddress: &emailpbv1.SenderAccount{
						Realm: "",
						Name:  "",
						Email: "sender@mail.com",
					},
					addresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{
								Name:  "",
								Email: "some@mail.com",
							},
						},
					},
					subject: "subject",
					body:    "body",
					data:    nil,
				}
			},
			expErr: nil,
		},
		{
			name:     "err no subject",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) msgData {
				return msgData{
					fromAddress: &emailpbv1.SenderAccount{
						Realm: "",
						Name:  "",
						Email: "sender@mail.com",
					},
					addresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{
								Name:  "",
								Email: "some@mail.com",
							},
						},
					},
					subject: "",
					body:    "",
					data:    nil,
				}
			},
			expErr: fmt.Errorf("asd"),
		},
		{
			name:     "err no body",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) msgData {
				return msgData{
					fromAddress: &emailpbv1.SenderAccount{
						Realm: "",
						Name:  "",
						Email: "sender@mail.com",
					},
					addresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{
								Name:  "",
								Email: "some@mail.com",
							},
						},
					},
					subject: "some subject",
					body:    "",
					data:    nil,
				}
			},
			expErr: fmt.Errorf("asd"),
		},
		{
			name:     "err no from address",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) msgData {
				return msgData{
					addresses: &emailpbv1.ToAddresses{
						ToEmail: []*emailpbv1.EmailAddr{
							{
								Name:  "",
								Email: "some@mail.com",
							},
						},
					},
					subject: "some subject",
					body:    "body",
					data:    nil,
				}
			},
			expErr: fmt.Errorf("asd"),
		},
		{
			name:     "err no to address",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) msgData {
				return msgData{
					fromAddress: &emailpbv1.SenderAccount{},
					subject:     "some subject",
					body:        "body",
					data:        nil,
				}
			},
			expErr: fmt.Errorf("asd"),
		},
		{
			name:     "err empty from address",
			expected: "to:some@mail.com\nfrom:sender@mail.com\nheader:From:sender@mail.com\nheader:To:some@mail.com\nheader:Subject:subject\nbody",
			from: func(t *testing.T) msgData {
				return msgData{
					fromAddress: &emailpbv1.SenderAccount{},
					addresses:   &emailpbv1.ToAddresses{},
					subject:     "some subject",
					body:        "body",
					data:        nil,
				}
			},
			expErr: fmt.Errorf("asd"),
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
				require.Error(t, err)
				fmt.Printf("%s %v %#v \n", tc.name, err, err)
				fmt.Printf("%s %v %#v \n", tc.name, &tc.expErr, &tc.expErr)
			}
		})
	}

}
