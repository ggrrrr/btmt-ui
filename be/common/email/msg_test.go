package email

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			name:    "ok",
			from:    Rcpt{Mail: "from@me", Name: "c"},
			to:      RcptList{Rcpt{Mail: "to@me", Name: "to"}, {Mail: "to1@me", Name: "to1"}},
			subject: "subject1",
			want: &Msg{
				from: Rcpt{Mail: "from@me", Name: "c"},
				to:   RcptList{Rcpt{Mail: "to@me", Name: "to"}, {Mail: "to1@me", Name: "to1"}},
				headers: []smtpHeader{
					smtpHeader{key: headerFrom, values: []string{"\"c\" <from@me>"}},
					smtpHeader{key: headerTo, values: []string{"\"to\" <to@me>", "\"to1\" <to1@me>"}},
					smtpHeader{key: headerSubject, values: []string{"subject1"}},
				},
				parts:       []*mailPart{},
				attachments: []*attachmentPart{},
				encoding:    "quoted-printable",
				charset:     "UTF-8",
			},
		},
		{
			name:    "from missing",
			from:    Rcpt{Mail: "from@me", Name: "c"},
			subject: "subject1",
			want:    &Msg{},
			err:     fmt.Errorf(""),
		},
		{
			name:    "to missing",
			from:    Rcpt{Mail: "from@me", Name: "c"},
			subject: "",
			want:    &Msg{},
			err:     fmt.Errorf(""),
		},
		{
			name:    "subject missing",
			from:    Rcpt{Mail: "from@me", Name: "c"},
			to:      RcptList{Rcpt{Mail: "to@me", Name: "to"}, {Mail: "to1@me", Name: "to1"}},
			subject: "",
			want:    &Msg{},
			err:     fmt.Errorf(""),
		},
	}
	for _, tc := range tst {
		t.Run(tc.name, func(t *testing.T) {
			m, gotErr := CreateMsg(tc.from, tc.to, tc.subject)
			if tc.err != nil {
				assert.Error(t, gotErr)
			} else {
				require.NoError(t, gotErr)
				assert.Equal(t, tc.want, m)
			}
		})
	}
}
