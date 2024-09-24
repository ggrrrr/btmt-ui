package email

import "fmt"

type (
	Rcpt struct {
		// mail@to.com
		Mail string
		// Mail name
		Name string
	}

	RcptList []Rcpt
)

// return "Mail name" <email.com>
func (rcpt Rcpt) Format() string {
	if rcpt.Name == "" {
		return rcpt.Mail
	}
	return fmt.Sprintf("\"%s\" <%s>", rcpt.Name, rcpt.Mail)
}

// return mail1@root,mail2@root
func (r RcptList) JoinMails() string {
	if r == nil {
		return ""
	}
	if len(r) == 0 {
		return ""
	}
	out := ""
	for _, v := range r {
		out += fmt.Sprintf("%s,", v.Mail)
	}
	return out[:len(out)-1]
}

// return "Mail name" <email.com>,"Mail name 1" <email1.com>
func (r RcptList) FormatedMails() []string {
	if r == nil {
		return []string{}
	}
	if len(r) == 0 {
		return []string{}
	}
	out := []string{}
	for _, v := range r {
		out = append(out, v.Format())
	}
	return out
}

func (r RcptList) Format() []string {
	if r == nil {
		return []string{}
	}
	if len(r) == 0 {
		return []string{}
	}
	out := []string{}
	for _, v := range r {
		out = append(out, v.Format())
	}
	return out
}
