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

// Implements io.Writer interface
func (w *partWriter) Write(bytes []byte) (int, error) {
	return w.w.Write(bytes)
}

func (w *partWriter) writeString(str string) error {
	_, err := w.w.Write([]byte(str))
	return err
}

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

// TODO implment quotedprintable
// Implmement multiple parts
func (w *partWriter) writePart(part *mailPart) error {
	w.writeBoundaryStart()
	w.writeHeader(headerContentType, string(part.contentType))
	w.Write(newLine)

	// qp := quotedprintable.NewWriter(w.w)
	// defer qp.Close()
	return part.copier(w.w)
}

// TODO implment text files none base64 encoding based on mimeType
func (w *partWriter) writeAttachment(part *attachment) error {
	mediaType := mime.TypeByExtension(filepath.Ext(part.name))

	w.writeBoundaryStart()
	w.writeHeader(headerContentType, fmt.Sprintf(`%s; name="%s"`, mediaType, part.name))
	w.writeHeader(headerContentTransferEncoding, string(Base64))
	w.writeHeader(headerContentDisposition, fmt.Sprintf(`attachment; filename="%s"`, part.name))
	w.writeHeader(headerContentID, fmt.Sprintf(`<%s>`, part.name))
	w.Write(newLine)

	wc := base64.NewEncoder(base64.StdEncoding, w)
	defer wc.Close()
	return part.copier(wc)
}
