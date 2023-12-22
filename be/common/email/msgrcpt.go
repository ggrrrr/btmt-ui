package email

import "fmt"

type (
	Rcpt struct {
		Mail string
		Name string
	}

	RcptList []Rcpt
)

func (rcpt Rcpt) Format() string {
	if rcpt.Name == "" {
		return rcpt.Mail
	}
	return fmt.Sprintf("\"%s\" <%s>", rcpt.Name, rcpt.Mail)
}

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

func (r RcptList) Mails() []string {
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
