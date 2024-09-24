package email

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"path/filepath"
)

type (
	partWriter struct {
		w        io.Writer
		boundary string
	}
)

var (
	newLine []byte = []byte("\r\n")
)

func (w *partWriter) writeBoundaryStart() error {
	return w.writeString(fmt.Sprintf("%s--%s%s", newLine, w.boundary, newLine))
}

func (w *partWriter) writeBoundaryClose() error {
	return w.writeString(fmt.Sprintf("%s--%s--%s", newLine, w.boundary, newLine))
}

func (w *partWriter) writeHeader(header headerName, values ...string) error {
	var err error
	err = w.writeString(string(header))
	if err != nil {
		return err
	}
	err = w.writeString(": ")
	if err != nil {
		return err
	}
	for i, v := range values {
		if i > 0 {
			err = w.writeString(", ")
			if err != nil {
				return err
			}
		}
		err = w.writeString(v)
		if err != nil {
			return err
		}
	}
	_, err = w.Write(newLine)
	return err
}

// TODO implement quotedprintable
// Implement multiple parts
func (w *partWriter) writePart(part *mailPart) error {
	var err error
	err = w.writeBoundaryStart()
	if err != nil {
		return err
	}

	err = w.writeHeader(headerContentType, string(part.contentType))
	if err != nil {
		return err
	}

	_, err = w.Write(newLine)
	if err != nil {
		return err
	}

	// qp := quotedprintable.NewWriter(w.w)
	// defer qp.Close()
	return part.copier(w.w)
}

// TODO implement text files none base64 encoding based on mimeType
func (w *partWriter) writeAttachment(part *attachment) error {
	var err error
	mediaType := mime.TypeByExtension(filepath.Ext(part.name))

	err = w.writeBoundaryStart()
	if err != nil {
		return err
	}
	err = w.writeHeader(headerContentType, fmt.Sprintf(`%s; name="%s"`, mediaType, part.name))
	if err != nil {
		return err
	}

	err = w.writeHeader(headerContentTransferEncoding, string(Base64))
	if err != nil {
		return err
	}

	err = w.writeHeader(headerContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, part.name))
	if err != nil {
		return err
	}

	err = w.writeHeader(headerContentID, fmt.Sprintf(`<%s>`, part.name))
	if err != nil {
		return err
	}

	_, err = w.Write(newLine)
	if err != nil {
		return err
	}

	wc := base64.NewEncoder(base64.StdEncoding, w)
	defer wc.Close()
	return part.copier(wc)
}

func (w *partWriter) writeMimeVer10() error {
	return w.writeHeader(headerMimeVer, "1.0")
}

func (w *partWriter) writeMultipart(multipart string) error {
	return w.writeHeader(headerContentType, fmt.Sprintf("%s; boundary=%s", multipart, w.boundary))
}

// Implements io.Writer interface
func (w *partWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *partWriter) writeString(str string) error {
	_, err := w.w.Write([]byte(str))
	return err
}
