package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
)

type smtpClientMock struct {
	errorOnStartTLS  bool
	falseOnExtension bool
	errorOnAuth      bool
	buf              *bytes.Buffer
	dataBlocks       []string
	from             string
	to               []string
}

func (d *smtpClientMock) Write(p []byte) (n int, err error) {
	// fmt.Printf("data write-> %s||\n", string(p))
	return d.buf.Write(p)
}

func (d *smtpClientMock) Close() (err error) {
	if d.buf != nil {
		d.dataBlocks = append(d.dataBlocks, d.buf.String())
		d.buf = nil
	}
	return nil
}

var _ (extSmtpClient) = (*smtpClientMock)(nil)

func newSmtpClientMock() *smtpClientMock {
	return &smtpClientMock{
		buf:        nil,
		dataBlocks: []string{},
		to:         []string{},
	}
}

func (*smtpClientMock) Hello(string) error {
	return nil
}

func (s *smtpClientMock) Extension(string) (bool, string) {
	return true, ""
}

func (s *smtpClientMock) StartTLS(*tls.Config) error {
	if s.errorOnStartTLS {
		return fmt.Errorf("starttls")
	}
	return nil
}

func (s *smtpClientMock) Auth(smtp.Auth) error {
	if s.errorOnAuth {
		return fmt.Errorf("auth")
	}
	return nil
}

func (s *smtpClientMock) Mail(from string) error {
	s.from = from
	return nil
}

func (s *smtpClientMock) Rcpt(rcpt string) error {
	s.to = append(s.to, rcpt)
	return nil
}

func (d *smtpClientMock) Data() (io.WriteCloser, error) {
	d.buf = bytes.NewBufferString("")
	return d, nil
}

func (*smtpClientMock) Quit() error {
	return nil
}
