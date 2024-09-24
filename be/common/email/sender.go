package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type (
	AuthType string

	extSmtpClient interface {
		Hello(string) error
		Extension(string) (bool, string)
		StartTLS(*tls.Config) error
		Auth(smtp.Auth) error
		Mail(string) error
		Rcpt(string) error
		Data() (io.WriteCloser, error)
		Quit() error
		Close() error
	}

	Config struct {
		SMTPHost string
		SMTPAddr string
		Username string
		Password string
		AuthType AuthType
		Timeout  time.Duration
	}

	// Implements SenderCloser
	Sender struct {
		cfg        Config
		tcpConn    net.Conn
		smtpClient extSmtpClient
	}

	SenderCloser interface {
		Send(msg *Msg) error
		io.Closer
	}
)

const (
	AuthTypePlain AuthType = "plain"
)

var _ (SenderCloser) = (*Sender)(nil)

func NewSender(cfg Config) (*Sender, error) {
	if cfg.AuthType == "" {
		cfg.AuthType = AuthTypePlain
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	smtpConn := &Sender{
		cfg: cfg,
	}

	tcpConn, err := net.DialTimeout("tcp", cfg.SMTPAddr, cfg.Timeout)
	if err != nil {
		logger.Error(err).Int("timeoutSeconds", int(cfg.Timeout.Seconds())).Str("smtp_addr", cfg.SMTPAddr).Msg("net.Dial")
		return nil, fmt.Errorf("net.DialTimeout: %w", err)
	}
	smtpConn.tcpConn = tcpConn
	return smtpConn.dialAndAuth()
}

func (sender *Sender) dialAndAuth() (*Sender, error) {

	logger.Info().Str("smtp_addr", sender.cfg.SMTPAddr).Int("timeoutSeconds", int(sender.cfg.Timeout.Seconds())).Msg("connected")
	smtpClient, err := smtp.NewClient(sender.tcpConn, sender.cfg.SMTPHost)
	if err != nil {
		return nil, fmt.Errorf("smtp.NewClient: %w", err)
	}

	if ok, tlsInfo := smtpClient.Extension("STARTTLS"); ok {
		if err := smtpClient.StartTLS(&tls.Config{ServerName: sender.cfg.SMTPHost}); err != nil {
			smtpClient.Close()
			return nil, fmt.Errorf("smtpClient.StartTLS: %w", err)
		}
		logger.Info().Str("tls.info", tlsInfo).Msg("client.StartTLS")
	}

	auth := smtp.PlainAuth("", sender.cfg.Username, sender.cfg.Password, sender.cfg.SMTPHost)
	err = smtpClient.Auth(auth)
	if err != nil {
		logger.Error(err).Str("username", sender.cfg.Username).Msg("PlainAuth")
		return nil, fmt.Errorf("smtpClient.Auth: %w", err)
	}

	sender.smtpClient = smtpClient
	return sender, nil
}

func (a *Sender) Send(email *Msg) error {
	var err error
	err = a.smtpClient.Mail(email.from.Mail)
	if err != nil {
		return fmt.Errorf("smtpClient.Mail[%s]: %w", email.from.Mail, err)
	}

	for _, t := range email.to {
		if err := a.smtpClient.Rcpt(t.Mail); err != nil {
			return fmt.Errorf("smtpClient.to[].Rcpt[%s]: %w", t.Mail, err)
		}
	}

	w, err := a.smtpClient.Data()
	if err != nil {
		return fmt.Errorf("smtpClient.Data: %w", err)
	}
	defer w.Close()
	err = email.writerTo(w)
	if err != nil {
		return fmt.Errorf("email.writeTo[%s]: %w", email.to[0].Mail, err)
	}
	logger.Info().Str("to", email.to[0].Mail).Msg("Send")
	return nil
}

func (conn *Sender) Close() error {
	if err := conn.smtpClient.Quit(); err != nil {
		logger.Error(err).Msg("Close")
	}
	return conn.smtpClient.Close()
}
