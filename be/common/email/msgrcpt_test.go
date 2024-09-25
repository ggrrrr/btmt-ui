package email

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRcptCreate(t *testing.T) {
	tests := []struct {
		name     string
		fromAddr string
		expRcpt  Rcpt
		expErr   error
	}{
		{
			name:     "ok too many emails",
			fromAddr: "\"me\"",
			expRcpt:  Rcpt{addr: "me@mail.com", name: ""},
			expErr:   &MailFormatError{},
		},
		{
			name:     "ok no email",
			fromAddr: "\"me@asd.com,you@asd.com\"",
			expRcpt:  Rcpt{addr: "me@mail.com", name: ""},
			expErr:   &MailFormatError{},
		},
		{
			name:     "ok no name",
			fromAddr: "me@mail.com",
			expRcpt:  Rcpt{addr: "me@mail.com", name: ""},
			expErr:   nil,
		},
		{
			name:     "ok",
			fromAddr: "\"me\" <me@mail.com>",
			expRcpt:  Rcpt{addr: "me@mail.com", name: "me"},
			expErr:   nil,
		},
		{
			name:     "ok",
			fromAddr: "me@mail.com",
			expRcpt:  Rcpt{addr: "me@mail.com", name: ""},
			expErr:   nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RcptFromString(tc.fromAddr)
			if tc.expErr != nil {
				fmt.Printf("%#v %v\n", err, err)
				assert.Error(t, err)
				assert.ErrorAs(t, err, &tc.expErr)

			} else {
				assert.Equal(t, tc.expRcpt, got)
			}
		})
	}

}

func TestRcptList(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "rcpt format",
			testFunc: func(t *testing.T) {
				rcpt := Rcpt{addr: "me@com", name: "me"}
				res := rcpt.Formatted()
				assert.Equal(t, `"me" <me@com>`, res)
			},
		},
		{
			name: "rcpt format no name",
			testFunc: func(t *testing.T) {
				rcpt := Rcpt{addr: "me@com", name: ""}
				res := rcpt.Formatted()
				assert.Equal(t, `me@com`, res)
			},
		},
		{
			name: "rcptlist  create",
			testFunc: func(t *testing.T) {
				rcptList, err := RcptListFromString([]string{`"mail name" <mail@gmail.com>`, `mail@com`})
				require.NoError(t, err)
				assert.Equal(t,
					RcptList{
						Rcpt{addr: "mail@gmail.com", name: "mail name"},
						Rcpt{addr: "mail@com", name: ""},
					},
					rcptList)
				assert.Equal(t, rcptList.Formatted(), []string{"\"mail name\" <mail@gmail.com>", "mail@com"})
			},
		},
		{
			name: "rcptlist  create err",
			testFunc: func(t *testing.T) {
				_, err := RcptListFromString([]string{`"mail name" <mail@gmail.com>`, `<mail@com`})
				require.Error(t, err)
				fmt.Printf("%#v \n ", err)
			},
		},
		{
			name: "rcpt list format",
			testFunc: func(t *testing.T) {
				rcptList := RcptList{
					Rcpt{addr: "me1@com", name: "me 1"},
					Rcpt{addr: "me2@com", name: "me 2"},
				}
				resList := rcptList.Formatted()
				assert.Equal(t, []string{`"me 1" <me1@com>`, `"me 2" <me2@com>`}, resList)
				resAddr := rcptList.AddressList()
				assert.Equal(t, `me1@com,me2@com`, resAddr)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}

}
