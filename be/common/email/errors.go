package email

import "fmt"

type SmtpAuthError struct {
	host string
	err  error
	msg  string
}

func (a *SmtpAuthError) Error() string {
	return fmt.Sprintf("unable to %s host: %s err: %v", a.msg, a.host, a.err)
}

func (a *SmtpAuthError) Unwrap() error {
	return a.err
}

type SmtpConnError struct {
	host string
	err  error
	msg  string
}

func (a *SmtpConnError) Error() string {
	return fmt.Sprintf("unable to %s host: %s err: %v", a.msg, a.host, a.err)
}

func (a *SmtpConnError) Unwrap() error {
	return a.err
}
