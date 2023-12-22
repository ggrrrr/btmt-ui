package email

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type (
	Msg struct {
		from        Rcpt
		to          RcptList
		headers     headers
		parts       []*mailPart
		attachments []*attachment
		encoding    encoding
		charset     string
	}

	headers map[headerName][]string

	mailPart struct {
		contentType contentType
		copier      func(io.Writer) error
		encoding    encoding
	}

	attachment struct {
		name   string
		header map[string][]string
		copier func(w io.Writer) error
	}
)

func CreateMsg(from Rcpt, to RcptList, subject string) (*Msg, error) {
	if from.Mail == "" {
		return nil, fmt.Errorf("from is empty")
	}
	if len(to) == 0 {
		return nil, fmt.Errorf("to list is empty")
	}
	if to[0].Mail == "" {
		return nil, fmt.Errorf("to[0] is empty")
	}
	if subject == "" {
		return nil, fmt.Errorf("subject is empty")
	}
	msg := &Msg{
		from:        from,
		to:          to,
		headers:     headers{},
		parts:       []*mailPart{},
		attachments: []*attachment{},
		charset:     "UTF-8",
		encoding:    QuotedPrintable,
	}

	msg.setHeader(headerFrom, from.Format())
	msg.setHeader(headerTo, to.Mails()...)
	//return "=?utf-8?q?" + subject + "?="
	msg.setHeader(headerSubject, subject)

	return msg, nil
}

func (e *Msg) AddCc(cc RcptList) {
	e.to = append(e.to, cc...)
	e.setHeader(headerCc, cc.Format()...)
}

func (e *Msg) AddBcc(bcc RcptList) {
	e.to = append(e.to, bcc...)
}

func (e *Msg) AddBodyString(body string) {
	e.parts = append(e.parts, &mailPart{
		contentType: contentTypePlain,
		copier:      newStringCopier(body),
		encoding:    Unencoded,
	})
}

func (e *Msg) AddHtmlBodyString(body string) {
	e.parts = append(e.parts, &mailPart{
		contentType: contentTypeHtml,
		copier:      newStringCopier(body),
		encoding:    Unencoded,
	})
}

func (e *Msg) AddHtmlBodyWriter(copier func(io.Writer) error) {
	e.parts = append(e.parts, &mailPart{
		contentType: contentTypeHtml,
		copier:      copier,
		encoding:    Unencoded,
	})
}

func (e *Msg) AddAttachment(name string, copier func(io.Writer) error) {
	f := &attachment{
		name:   name,
		header: map[string][]string{},
		copier: copier,
	}
	e.attachments = append(e.attachments, f)
}

func (e *Msg) AddFile(fileName string) {
	f := &attachment{
		name:   filepath.Base(fileName),
		header: map[string][]string{},
		copier: func(w io.Writer) error {
			h, err := os.Open(fileName)
			if err != nil {
				return err
			}
			if _, err := io.Copy(w, h); err != nil {
				h.Close()
				return err
			}
			return h.Close()
		},
	}
	e.attachments = append(e.attachments, f)
}

func (m *Msg) setHeader(field headerName, value ...string) {
	m.headers[field] = value
}

func newStringCopier(s string) func(io.Writer) error {
	return func(w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	}
}
