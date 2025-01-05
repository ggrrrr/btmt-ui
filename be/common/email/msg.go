package email

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type (
	Msg struct {
		from        Rcpt
		to          RcptList
		headers     []smtpHeader
		parts       []*bodyPart
		attachments []*attachmentPart
		encoding    encoding
		charset     string
		rootWriter  *partWriter
	}

	smtpHeader struct {
		key    headerName
		values []string
	}

	bodyPart struct {
		contentType contentType
		copier      func(io.Writer) error
		encoding    encoding
	}

	attachmentPart struct {
		name   string
		header map[string][]string
		copier func(w io.Writer) error
	}
)

func CreateMsgFromString(from string, toList []string) (*Msg, error) {
	rcptFrom, err := RcptFromString(from)
	if err != nil {
		return nil, fmt.Errorf("from address: %w", err)
	}
	rcptTo := RcptList{}
	for _, v := range toList {
		rcpt, err := RcptFromString(v)
		if err != nil {
			return nil, fmt.Errorf("to address: %w", err)
		}
		rcptTo = append(rcptTo, rcpt)
	}

	return createMsg(rcptFrom, rcptTo)
}

func (msg *Msg) SetSubject(subject string) error {
	if subject == "" {
		return fmt.Errorf("subject is empty")
	}
	msg.setHeader(headerSubject, subject)
	return nil
}

func (e *Msg) AddCc(cc RcptList) {
	e.to = append(e.to, cc...)
	e.setHeader(headerCc, cc.Formatted()...)
}

func (e *Msg) AddBcc(bcc RcptList) {
	e.to = append(e.to, bcc...)
}

func (e *Msg) AddBodyString(body string) {
	e.parts = append(e.parts, &bodyPart{
		contentType: contentTypePlain,
		copier:      newStringCopier(body),
		encoding:    Unencoded,
	})
}

func (e *Msg) AddHtmlBodyString(body string) {
	e.parts = append(e.parts, &bodyPart{
		contentType: contentTypeHtml,
		copier:      newStringCopier(body),
		encoding:    Unencoded,
	})
}

func (e *Msg) AddHtmlBodyWriter(copier func(io.Writer) error) {
	e.parts = append(e.parts, &bodyPart{
		contentType: contentTypeHtml,
		copier:      copier,
		encoding:    Unencoded,
	})
}

func (e *Msg) AddAttachment(name string, copier func(io.Writer) error) {
	f := &attachmentPart{
		name:   name,
		header: map[string][]string{},
		copier: copier,
	}
	e.attachments = append(e.attachments, f)
}

func (e *Msg) AddFile(fileName string) {
	f := &attachmentPart{
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

func createMsg(from Rcpt, to RcptList) (*Msg, error) {
	if from.addr == "" {
		return nil, fmt.Errorf("from is empty")
	}
	if to == nil {
		return nil, fmt.Errorf("to list is nil")
	}
	if len(to) == 0 {
		return nil, fmt.Errorf("to list is empty")
	}
	if to[0].addr == "" {
		return nil, fmt.Errorf("to[0] is empty")
	}
	msg := &Msg{
		from:        from,
		to:          to,
		headers:     []smtpHeader{},
		parts:       []*bodyPart{},
		attachments: []*attachmentPart{},
		charset:     "UTF-8",
		encoding:    QuotedPrintable,
	}

	msg.setHeader(headerFrom, from.Formatted())
	msg.setHeader(headerTo, to.Formatted()...)

	return msg, nil
}

func (m *Msg) setHeader(field headerName, value ...string) {
	m.headers = append(m.headers, smtpHeader{
		key:    field,
		values: value,
	})
}

func newStringCopier(s string) func(io.Writer) error {
	return func(w io.Writer) error {
		_, err := io.WriteString(w, s)
		return err
	}
}

func (m *Msg) DumpToText() string {
	var buffer bytes.Buffer
	buffer.WriteString("to:")
	buffer.WriteString(m.to.AddressList())
	buffer.WriteString("\n")

	buffer.WriteString("from:")
	buffer.WriteString(m.from.Formatted())
	buffer.WriteString("\n")

	for _, v := range m.headers {
		buffer.WriteString("header:")
		buffer.WriteString(string(v.key))
		buffer.WriteString(":")
		buffer.WriteString(strings.Join(v.values, ","))
		buffer.WriteString("\n")
	}
	// buffer.WriteString("from:")
	// buffer.WriteString(m.headers)
	// buffer.WriteString("\n")

	for _, v := range m.parts {
		v.copier(&buffer)
	}

	// m.parts

	return buffer.String()
}
