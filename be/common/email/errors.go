package email

import "fmt"

type (
	SmtpAuthError struct {
		host string
		err  error
		msg  string
	}

	SmtpConnError struct {
		host string
		err  error
		msg  string
	}

	MailFormatError struct {
		msg string
		err error
	}
)

func (a *SmtpAuthError) Error() string {
	return fmt.Sprintf("unable to %s host: %s err: %v", a.msg, a.host, a.err)
}

func (a *SmtpAuthError) Unwrap() error {
	return a.err
}

func (a *SmtpConnError) Error() string {
	return fmt.Sprintf("unable to %s host: %s err: %v", a.msg, a.host, a.err)
}

func (a *SmtpConnError) Unwrap() error {
	return a.err
}

func (a *MailFormatError) Error() string {
	return fmt.Sprintf("mail %s: %v", a.msg, a.err)
}

func (a *MailFormatError) Unwrap() error {
	return a.err
}
