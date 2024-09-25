package email

import (
	"fmt"
	"net/mail"
)

type (
	Rcpt struct {
		// mail@to.com
		addr string
		// Mail name
		name string
	}

	RcptList []Rcpt
)

func RcptFromString(addr string) (Rcpt, error) {
	fromAddr, err := mail.ParseAddress(addr)
	if err != nil {
		return Rcpt{}, &MailFormatError{
			msg: "email address",
			err: err,
		}
	}
	return Rcpt{
		addr: fromAddr.Address,
		name: fromAddr.Name,
	}, nil

}

func RcptListFromString(addr []string) (list RcptList, err error) {
	for _, v := range addr {
		rcpt, err := RcptFromString(v)
		if err != nil {
			return nil, fmt.Errorf("rcpt list: %w", err)
		}
		list = append(list, rcpt)

	}
	return

}

// return "Mail name" <email.com>
func (rcpt Rcpt) Formatted() string {
	if rcpt.name == "" {
		return rcpt.addr
	}
	return fmt.Sprintf("\"%s\" <%s>", rcpt.name, rcpt.addr)
}

// return mail1@root,mail2@root
func (r RcptList) AddressList() string {
	if r == nil {
		return ""
	}
	if len(r) == 0 {
		return ""
	}
	out := ""
	for _, v := range r {
		out += fmt.Sprintf("%s,", v.addr)
	}
	return out[:len(out)-1]
}

// return []{"\"Mail name\" <email.com>","\"Mail name 1\" <email1.com>""}
func (r RcptList) Formatted() []string {
	if r == nil {
		return []string{}
	}
	if len(r) == 0 {
		return []string{}
	}
	out := []string{}
	for _, v := range r {
		out = append(out, v.Formatted())
	}
	return out
}
