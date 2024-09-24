package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/smtp"
)

type SmtpClientMock struct {
	errorOnStartTLS  bool
	falseOnExtension bool
	errorOnAuth      bool
	buf              *bytes.Buffer
	dataBlocks       []string
}

func (d *SmtpClientMock) Write(p []byte) (n int, err error) {
	// fmt.Printf("data write-> %s||\n", string(p))
	return d.buf.Write(p)
}

func (d *SmtpClientMock) Close() (err error) {
	if d.buf != nil {
		d.dataBlocks = append(d.dataBlocks, d.buf.String())
		d.buf = nil
	}
	return nil
}

var _ (extSmtpClient) = (*SmtpClientMock)(nil)

func NewSmtpClientMock() *SmtpClientMock {
	return &SmtpClientMock{
		buf:        nil,
		dataBlocks: []string{},
	}
}

func (*SmtpClientMock) Hello(string) error {
	return nil
}

func (s *SmtpClientMock) Extension(string) (bool, string) {
	return true, ""
}

func (s *SmtpClientMock) StartTLS(*tls.Config) error {
	if s.errorOnStartTLS {
		return fmt.Errorf("starttls")
	}
	return nil
}

func (s *SmtpClientMock) Auth(smtp.Auth) error {
	if s.errorOnAuth {
		return fmt.Errorf("auth")
	}
	return nil
}

func (*SmtpClientMock) Mail(string) error {
	return nil
}

func (*SmtpClientMock) Rcpt(string) error {
	return nil
}

func (d *SmtpClientMock) Data() (io.WriteCloser, error) {
	d.buf = bytes.NewBufferString("")
	return d, nil
}

func (*SmtpClientMock) Quit() error {
	return nil
}
