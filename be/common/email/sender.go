package email

import (
	"crypto/tls"
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
	sender := &Sender{
		cfg: cfg,
	}

	tcpConn, err := net.DialTimeout("tcp", cfg.SMTPAddr, cfg.Timeout)
	if err != nil {
		logger.Error(err).Int("timeoutSeconds", int(cfg.Timeout.Seconds())).Str("smtp_addr", cfg.SMTPAddr).Msg("net.Dial")
		return nil, &SmtpConnError{
			host: cfg.SMTPHost,
			err:  err,
			msg:  "Dial",
		}
	}
	sender.tcpConn = tcpConn

	logger.Info().Str("smtp_addr", sender.cfg.SMTPAddr).Int("timeoutSeconds", int(sender.cfg.Timeout.Seconds())).Msg("connected")
	smtpClient, err := smtp.NewClient(sender.tcpConn, sender.cfg.SMTPHost)
	if err != nil {
		// return nil, fmt.Errorf("smtp.NewClient: %w", err)
		return nil, &SmtpConnError{
			host: cfg.SMTPHost,
			err:  err,
			msg:  "unable create smtp client",
		}
	}

	err = sender.smtpAuth()
	if err != nil {
		smtpClient.Quit()
		return nil, err
	}
	sender.smtpClient = smtpClient
	return sender, err
}

func (sender *Sender) smtpAuth() error {
	var err error

	ok, tlsInfo := sender.smtpClient.Extension("STARTTLS")
	if ok {
		err = sender.smtpClient.StartTLS(&tls.Config{ServerName: sender.cfg.SMTPHost})
		if err != nil {
			sender.smtpClient.Close()
			return &SmtpAuthError{
				host: sender.cfg.SMTPHost,
				err:  err,
				msg:  "STARTTLS",
			}
		}
		logger.Info().Str("tls.info", tlsInfo).Msg("client.StartTLS")
	} else {
		logger.Info().
			Str("tls.info", tlsInfo).
			Bool("STARTTLS", ok).
			Msg("client.StartTLS")
	}

	auth := smtp.PlainAuth("", sender.cfg.Username, sender.cfg.Password, sender.cfg.SMTPHost)
	err = sender.smtpClient.Auth(auth)
	if err != nil {
		logger.Error(err).Str("username", sender.cfg.Username).Msg("PlainAuth")
		return &SmtpAuthError{
			host: sender.cfg.SMTPHost,
			err:  err,
			msg:  "Auth",
		}
	}

	return nil
}

func (conn *Sender) Close() error {
	if err := conn.smtpClient.Quit(); err != nil {
		logger.Error(err).Msg("Close")
	}
	return conn.smtpClient.Close()
}