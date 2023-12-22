package email

import (
	"crypto/tls"
	"io"
	"net"
	"net/smtp"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
	"github.com/stackus/errors"
)

type (
	AuthType string

	Config struct {
		Host     string
		Addr     string
		Username string
		Password string
		AuthType AuthType
		timeout  time.Duration
	}

	// Implements SenderCloser
	smtpConn struct {
		cfg    Config
		client extSmtpClient
	}

	SenderCloser interface {
		Send(msg *Msg) error
		io.Closer
	}
)

const (
	AuthTypePlain AuthType = "plain"
)

var _ (SenderCloser) = (*smtpConn)(nil)

func Dial(cfg Config) (*smtpConn, error) {
	if cfg.AuthType == "" {
		cfg.AuthType = AuthTypePlain
	}
	if cfg.timeout == 0 {
		cfg.timeout = 10 * time.Second
	}
	smtpConn := &smtpConn{
		cfg: cfg,
	}
	return smtpConn.Dial()
}

func (smtpConn *smtpConn) Dial() (*smtpConn, error) {

	tcpConn, err := net.DialTimeout("tcp", smtpConn.cfg.Addr, smtpConn.cfg.timeout)
	if err != nil {
		return nil, errors.Wrap(err, "DialTimeout")
	}
	logger.Info().Str("host", smtpConn.cfg.Addr).Msg("Connected.")
	client, err := smtp.NewClient(tcpConn, smtpConn.cfg.Host)
	if err != nil {
		return nil, errors.Wrap(err, "NewClient")
	}

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: smtpConn.cfg.Host}); err != nil {
			client.Close()
			return nil, errors.Wrap(err, "StartTLS")
		}
		logger.Debug().Msg("StartTLS.")
	}

	auth := smtp.PlainAuth("", smtpConn.cfg.Username, smtpConn.cfg.Password, smtpConn.cfg.Host)
	err = client.Auth(auth)
	if err != nil {
		logger.Error(err).Str("username", smtpConn.cfg.Username).Msg("PlainAuth")
		return nil, errors.Wrap(err, "PlainAuth")
	}

	smtpConn.client = client
	return smtpConn, nil
}

func (a *smtpConn) Send(email *Msg) error {
	var err error
	err = a.client.Mail(email.from.Mail)
	if err != nil {
		return errors.Wrap(err, "Mail")
	}

	for _, t := range email.to {
		if err := a.client.Rcpt(t.Mail); err != nil {
			return errors.Wrap(err, "Rcpt")
		}
	}

	w, err := a.client.Data()
	if err != nil {
		return errors.Wrap(err, "Data")
	}
	defer w.Close()
	err = email.writerTo(w)
	if err != nil {
		return errors.Wrap(err, "WriterTo")
	}
	logger.Info().Str("to", email.to[0].Mail).Msg("Send")
	return nil
}

func (conn *smtpConn) Close() error {
	conn.client.Quit()
	return conn.client.Close()
}
