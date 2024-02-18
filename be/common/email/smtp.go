package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/smtp"
)

type (
	headerName  string
	encoding    string
	contentType string

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

func (e *Msg) writerTo(w io.Writer) error {
	if len(e.parts) == 0 {
		return fmt.Errorf("msg.parts is empty")
	}
	var err error
	rootPart := &partWriter{
		w: w,
	}
	mpWriter := multipart.NewWriter(rootPart.w)

	for k, v := range e.headers {
		err = rootPart.writeHeader(k, v...)
		if err != nil {
			return err
		}
	}
	err = rootPart.writeMimeVer10()
	if err != nil {
		return err
	}
	rootPart.boundary = mpWriter.Boundary()
	err = rootPart.writeMultipart(multipartRelated)
	if err != nil {
		return err
	}

	err = rootPart.writePart(e.parts[0])
	if err != nil {
		return err
	}

	for _, a := range e.attachments {
		err = rootPart.writeAttachment(a)
		if err != nil {
			return err
		}
	}
	return rootPart.writeBoundaryClose()
}

func (w *partWriter) writeMimeVer10() error {
	return w.writeHeader(headerMimeVer, "1.0")
}

func (w *partWriter) writeMultipart(multipart string) error {
	return w.writeHeader(headerContentType, fmt.Sprintf("%s; boundary=%s", multipart, w.boundary))
}
