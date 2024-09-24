package email

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/ggrrrr/btmt-ui/be/common/logger"
)

type (
	headerName  string
	encoding    string
	contentType string
)

const (
	multipartRelated     string = "multipart/related"
	multipartAlternative string = "multipart/alternative"
	multipartMixed       string = "multipart/mixed"

	contentTypeHtml  contentType = "text/html"
	contentTypePlain contentType = "text/plain"

	headerMimeVer                 headerName = "MIME-Version"
	headerSubject                 headerName = "Subject"
	headerTo                      headerName = "To"
	headerFrom                    headerName = "From"
	headerCc                      headerName = "Cc"
	headerContentType             headerName = "Content-Type"
	headerContentTransferEncoding headerName = "Content-Transfer-Encoding"
	headerContentDisposition      headerName = "Content-Disposition"
	headerContentID               headerName = "Content-ID"

	// QuotedPrintable represents the quoted-printable encoding as defined in
	// RFC 2045.
	QuotedPrintable encoding = "quoted-printable"
	// Base64 represents the base64 encoding as defined in RFC 2045.
	Base64 encoding = "base64"
	// Unencoded can be used to avoid encoding the body of an email. The headers
	// will still be encoded using quoted-printable encoding.
	Unencoded encoding = "8bit"
)

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

func (e *Msg) writerTo(w io.Writer) error {
	if len(e.parts) == 0 {
		return fmt.Errorf("msg.parts is empty")
	}
	var err error
	e.rootWriter = &partWriter{
		w:        w,
		boundary: randomBoundary(),
	}

	for _, v := range e.headers {
		err = e.rootWriter.writeHeader(v)
		if err != nil {
			return fmt.Errorf("msg.writeHeaders: %w", err)
		}
	}

	err = e.rootWriter.writeMimeVer10()
	if err != nil {
		return fmt.Errorf("msg.writeMimeVer:1.0: %w", err)
	}

	err = e.rootWriter.writeMultipart(multipartRelated)
	if err != nil {
		return fmt.Errorf("msg.writeMultipart: %w", err)
	}

	err = e.rootWriter.writePart(e.parts[0])
	if err != nil {
		return fmt.Errorf("msg.writePart[0]: %w", err)
	}

	for _, a := range e.attachments {
		err = e.rootWriter.writeAttachment(a)
		if err != nil {
			return fmt.Errorf("msg.writeAttachment[...]: %w", err)
		}
	}

	err = e.rootWriter.writeBoundaryClose()
	if err != nil {
		return fmt.Errorf("msg.writeBoundaryClose: %w", err)
	}

	return nil
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		logger.Error(err).Msg("randomBoundary ReadFull(rand.Reader)")
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}
