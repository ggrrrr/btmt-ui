package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/smtp"
	"time"

	"github.com/ggrrrr/btmt-ui/be/common/ltm/log"
	"github.com/ggrrrr/btmt-ui/be/common/ltm/tracer"
)

const otelScope string = "go.github.com.ggrrrr.btmt-ui.common.email"

const (
	AuthTypePlain AuthType = "plain"
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
		otelTracer tracer.OTelTracer
		cfg        Config
		tcpConn    net.Conn
		smtpClient extSmtpClient
	}

	EmailSender interface {
		Send(ctx context.Context, email *Msg) error
		Close() error
	}

	EmailConnector interface {
		Connect(ctx context.Context) (EmailSender, error)
	}
)

var _ (EmailConnector) = (*Config)(nil)

var _ (EmailSender) = (*Sender)(nil)

func (cfg *Config) Connect(ctx context.Context) (EmailSender, error) {
	var err error
	otelTracer := tracer.Tracer(otelScope)

	_, span := otelTracer.SpanWithAttributes(ctx, "email.NewSender",
		slog.String("smtp.host.addr", cfg.SMTPAddr),
	)

	defer func() {
		span.End(err)
	}()

	if cfg.AuthType == "" {
		cfg.AuthType = AuthTypePlain
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	sender := &Sender{
		cfg: *cfg,
	}

	tcpConn, err := net.DialTimeout("tcp", cfg.SMTPAddr, cfg.Timeout)
	if err != nil {
		log.Log().ErrorCtx(ctx, err, "net.Dial",
			log.WithString("smtp.addr", cfg.SMTPAddr),
			log.WithInt("smtp.timeout.sec", int(cfg.Timeout.Seconds())),
		)

		return nil, &SmtpConnError{
			host: cfg.SMTPHost,
			err:  err,
			msg:  "Dial",
		}
	}
	sender.tcpConn = tcpConn

	log.Log().InfoCtx(ctx, "net.Dial",
		log.WithString("smtp.addr", cfg.SMTPAddr),
		log.WithInt("smtp.timeout.sec", int(cfg.Timeout.Seconds())),
	)

	smtpClient, err := smtp.NewClient(sender.tcpConn, sender.cfg.SMTPHost)
	if err != nil {
		tcpConn.Close()
		sender.tcpConn = nil
		// return nil, fmt.Errorf("smtp.NewClient: %w", err)
		return nil, &SmtpConnError{
			host: cfg.SMTPHost,
			err:  err,
			msg:  "unable create smtp client",
		}
	}

	sender.smtpClient = smtpClient

	err = sender.smtpAuth()
	if err != nil {
		e := smtpClient.Quit()
		tcpConn.Close()
		sender.tcpConn = nil
		sender.smtpClient = nil

		if e != nil {
			log.Log().WarnCtx(ctx, err, "net.Connect",
				log.WithString("smtp.addr", cfg.SMTPAddr),
			)
		}
		return nil, err
	}
	return sender, err
}

func (a *Sender) Send(ctx context.Context, email *Msg) error {
	if len(email.parts) == 0 {
		return &MailFormatError{
			msg: "body",
			err: fmt.Errorf("is empty"),
		}
	}

	var err error
	_, span := a.otelTracer.SpanWithAttributes(ctx, "email.Send",
		slog.String("email.from", email.from.addr),
		slog.String("email.to", email.to.AddressList()),
	)

	defer func() {
		span.End(err)
	}()

	err = a.smtpClient.Mail(email.from.addr)
	if err != nil {
		return fmt.Errorf("smtpClient.Mail[%s]: %w", email.from.addr, err)
	}

	for _, t := range email.to {
		if err := a.smtpClient.Rcpt(t.addr); err != nil {
			return fmt.Errorf("smtpClient.to[].Rcpt[%s]: %w", t.addr, err)
		}
	}

	w, err := a.smtpClient.Data()
	if err != nil {
		return fmt.Errorf("smtpClient.Data: %w", err)
	}

	defer w.Close()
	err = email.writerTo(w)
	if err != nil {
		return fmt.Errorf("email.writeTo[%s]: %w", email.to[0].addr, err)
	}

	log.Log().InfoCtx(ctx, "Send",
		log.WithString("to", email.to[0].addr),
	)

	return nil
}

func (conn *Sender) Close() error {
	var err error
	if conn.smtpClient != nil {
		err = conn.smtpClient.Quit()
		conn.smtpClient = nil
	}
	if err != nil {
		log.Log().Error(err, "Close")
	}
	if conn.tcpConn != nil {
		err = conn.tcpConn.Close()
		conn.tcpConn = nil
	}
	if err != nil {
		log.Log().Error(err, "Close")
	}
	return nil
}

func (sender *Sender) smtpAuth() error {
	var err error

	log.Log().Debug("sender.smtpAuth")

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

		log.Log().Debug("sender.smtpAuth", log.WithAny("tls.info", tlsInfo))

	} else {
		log.Log().Debug("sender.client.StartTLS",
			log.WithAny("tls.info", tlsInfo),
			log.WithBool("tls.info.STARTTLS", ok),
		)
	}

	auth := smtp.PlainAuth("", sender.cfg.Username, sender.cfg.Password, sender.cfg.SMTPHost)
	err = sender.smtpClient.Auth(auth)
	if err != nil {

		log.Log().Debug("sender.client.PlainAuth",
			log.WithAny("username", sender.cfg.Username),
		)
		return &SmtpAuthError{
			host: sender.cfg.SMTPHost,
			err:  err,
			msg:  "Auth",
		}
	}

	log.Log().Debug("Auth.ok.")

	return nil
}
