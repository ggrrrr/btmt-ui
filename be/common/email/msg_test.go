package email

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDumpToText(t *testing.T) {

	msg := Msg{
		from: Rcpt{addr: "from@amil.com", name: "from"},
		to:   RcptList{Rcpt{addr: "to@mail.com", name: "to"}},
		headers: []smtpHeader{
			{key: "a", values: []string{"b"}},
		},
		parts:       []*bodyPart{{contentType: "", copier: newStringCopier("body text")}},
		attachments: []*attachmentPart{},
		encoding:    QuotedPrintable,
		charset:     "",
	}

	actual := msg.DumpToText()

	require.Equal(t, "to:to@mail.com\nfrom:\"from\" <from@amil.com>\nheader:a:b\nbody text", actual, "%v", actual)

}

func TestRcptMail(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "ok",
			testFunc: func(t *testing.T) {
				msg, err := CreateMsgFromString("me@me.com", []string{"to@to.com"})
				require.NoError(t, err)
				err = msg.SetSubject("subject")
				require.NoError(t, err)
				assert.Equal(t, &Msg{
					from: Rcpt{addr: "me@me.com"},
					to:   RcptList{Rcpt{addr: "to@to.com"}},
					headers: []smtpHeader{
						{key: "From", values: []string{"me@me.com"}},
						{key: "To", values: []string{"to@to.com"}},
						{key: "Subject", values: []string{"subject"}},
					},
					parts:       []*bodyPart{},
					attachments: []*attachmentPart{},
					encoding:    "quoted-printable",
					charset:     "UTF-8",
					rootWriter:  nil,
				}, msg)
			},
		},
		{
			name: "rcpt list fail",
			testFunc: func(t *testing.T) {
				_, err := CreateMsgFromString("me@me.com", []string{"to@to.com", "asdasd"})
				require.Error(t, err)
				err1 := &MailFormatError{}
				require.ErrorAs(t, err, &err1)
				fmt.Printf("%+v \n", err)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}

}

func TestCreateMsg(t *testing.T) {
	type testCase struct {
		name    string
		from    Rcpt
		to      RcptList
		subject string
		want    *Msg
		err     error
	}
	tst := []testCase{
		{
			name: "ok",
			from: Rcpt{addr: "from@me", name: "c"},
			to:   RcptList{Rcpt{addr: "to@me", name: "to"}, {addr: "to1@me", name: "to1"}},
			// subject: "subject1",
			want: &Msg{
				from: Rcpt{addr: "from@me", name: "c"},
				to:   RcptList{Rcpt{addr: "to@me", name: "to"}, {addr: "to1@me", name: "to1"}},
				headers: []smtpHeader{
					{key: headerFrom, values: []string{"\"c\" <from@me>"}},
					{key: headerTo, values: []string{"\"to\" <to@me>", "\"to1\" <to1@me>"}},
					// {key: headerSubject, values: []string{"subject1"}},
				},
				parts:       []*bodyPart{},
				attachments: []*attachmentPart{},
				encoding:    "quoted-printable",
				charset:     "UTF-8",
			},
		},
		{
			name:    "from missing",
			from:    Rcpt{addr: "from@me", name: "c"},
			subject: "subject1",
			want:    &Msg{},
			err:     fmt.Errorf(""),
		},
		{
			name:    "to missing",
			from:    Rcpt{addr: "from@me", name: "c"},
			subject: "",
			want:    &Msg{},
			err:     fmt.Errorf(""),
		},
	}

	for _, tc := range tst {
		t.Run(tc.name, func(t *testing.T) {
			m, gotErr := createMsg(tc.from, tc.to)
			// m, gotErr := createMsg(tc.from, tc.to, tc.subject)
			if tc.err != nil {
				assert.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
				assert.Equal(t, tc.want, m)
			}
		})
	}

}
